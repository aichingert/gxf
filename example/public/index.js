async function init() {
    const go = new window.Go();

    WebAssembly.instantiateStreaming(
        fetch("main.wasm"),
        go.importObject
    ).then(async (obj) => {
        go.run(obj.instance);

        const content = await fetch("test.dxf");
        const input = new Uint8Array(await content.arrayBuffer());
        const plan = window.parse(input);

        setupWGPU(plan);
    });
}

async function setupWGPU(plan) {
    let canvas = document.getElementById("drawing");
    canvas.width = canvas.clientWidth * window.devicePixelRatio;
    canvas.height = canvas.clientHeight * window.devicePixelRatio; 

    if (!navigator.gpu) {
        const errorHeader = document.createElement("h1");
        errorHeader.style.color = "white";

        const errorText = document.createTextNode(
            "This browser does not support Wgpu!"
        );

        errorHeader.appendChild(errorText);
        document.body.insertBefore(errorHeader, document.getElementById("dummy"));

        return;
    }

    const shaders_src = `
    struct Camera {
        denom: vec2f,
        small: vec2f,
    }

    struct VertexOut {
        @builtin(position) position : vec4f,
        @location(0) color : vec4f
    }

    @binding(0) @group(0) var<uniform> camera: Camera;

    @vertex
    fn vertex_main(@location(0) position: vec2f,
                   @location(1) color: vec3f) -> VertexOut
    {
        var x: f32 = (position.x - camera.small.x) / camera.denom.x - 1.0;
        var y: f32 = (position.y - camera.small.y) / camera.denom.y - 1.0;

        var output : VertexOut;
        output.position = vec4f(x, y, 1.0, 1.0);
        output.color = vec4f(color, 1.0);
        return output;
    }

    @fragment
    fn fragment_main(fragData: VertexOut) -> @location(0) vec4f
    {
        return fragData.color;
    }
    `; 

    let adapter;

    try {
        adapter = await navigator.gpu.requestAdapter();
    } catch (e) {
        const errorHeader = document.createElement("h1");
        errorHeader.style.color = "white";

        const errorText = document.createTextNode(
            "This browser does not support Wgpu!"
        );

        errorHeader.appendChild(errorText);
        document.body.insertBefore(errorHeader, document.getElementById("dummy"));

        return;
    }

    const device = await adapter.requestDevice();

    const shaderModule = device.createShaderModule({
        code: shaders_src,
    });

    let offset = 0;
    if (plan.BlockOffsets.length > 0) {
        offset = plan.BlockOffsets[plan.BlockOffsets.length - 1][0];
    }

    const full = new Float32Array(plan.Data.Lines.Vertices);
    const vertices = new Float32Array(full.subarray(offset, full.length));

    const vertexBuffer = device.createBuffer({
        size: vertices.byteLength,
        usage: GPUBufferUsage.VERTEX | GPUBufferUsage.COPY_DST,
    });

    device.queue.writeBuffer(vertexBuffer, 0, vertices, 0, vertices.length);

    const vertexBuffers = [
        {
            attributes: [
                {
                    shaderLocation: 0,
                    offset: 0,
                    format: "float32x2",
                },
                {
                    shaderLocation: 1,
                    offset: 8,
                    format: "float32x3",
                },
            ],
            arrayStride: 20,
            stepMode: "vertex",
        },
    ];

    const uniformBuffer = device.createBuffer({
        size: 16,
        usage: GPUBufferUsage.UNIFORM | GPUBufferUsage.COPY_DST,
    })

    const bindGroupLayout = device.createBindGroupLayout({
        entries: [
            {
                binding: 0,
                visibility: GPUShaderStage.VERTEX,
                buffer: {},
            }
        ],
    });

    const bindGroup = device.createBindGroup({
        layout: bindGroupLayout,
        entries: [
            {
                binding: 0,
                resource: {
                    buffer: uniformBuffer,
                }
            },
        ],
    });

    const pipelineLayout = device.createPipelineLayout({
        bindGroupLayouts: [bindGroupLayout],
    });

    const pipelineDescriptor = {
        vertex: {
            module: shaderModule,
            entryPoint: "vertex_main",
            buffers: vertexBuffers,
        },
        fragment: {
            module: shaderModule,
            entryPoint: "fragment_main",
            targets: [
                {
                    format: navigator.gpu.getPreferredCanvasFormat(),
                }
            ],
        },
        primitive: {
            topology: "line-list",
        },
        layout: pipelineLayout,
    };

    const renderPipeline = device.createRenderPipeline(pipelineDescriptor);

    const context = canvas.getContext("webgpu");
    context.configure({
        device: device,
        format: navigator.gpu.getPreferredCanvasFormat(),
        alphaMode: "premultiplied",
    });

    let s = 0.001;
    let fx = 1;

    console.log(plan);
    const xd = (plan.Data.MaxX - plan.Data.MinX) / 2.;
    const yd = (plan.Data.MaxY - plan.Data.MinY) / 2.;

    while (true) {
        await new Promise(r => setTimeout(r, 1));

        device.queue.writeBuffer(uniformBuffer, 0, new Float32Array([xd, yd, plan.Data.MinX, plan.Data.MinY]));

        const commandEncoder = device.createCommandEncoder();
        const clearColor = { r: 0.1289, g: 0.1289, b: 0.1289, a: 1.0 };

        const renderPassDescriptor = {
            colorAttachments: [
                {
                    clearValue: clearColor,
                    loadOp: "clear",
                    storeOp: "store",
                    view: context.getCurrentTexture().createView(),
                },
            ],
        };

        const passEncoder = commandEncoder.beginRenderPass(renderPassDescriptor);
        passEncoder.setPipeline(renderPipeline);
        passEncoder.setVertexBuffer(0, vertexBuffer);
        passEncoder.setBindGroup(0, bindGroup);
        passEncoder.draw(vertices.length / 5);

        passEncoder.end();
        device.queue.submit([commandEncoder.finish()]);
    }
}

window.onload = init;
