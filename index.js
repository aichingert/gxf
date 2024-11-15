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
        position: vec4f,
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
        var x : f32 = camera.position.x * position.x;
        var y : f32 = camera.position.y * position.y;

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

    const denY = (plan.MaxY - plan.MinY) / 2;
    const denX = (plan.MaxX - plan.MinX) / 2;

    console.log(plan);
    const lines = plan.Plan.Lines.Vertices;
    const mx = plan.Plan.MX;
    const my = plan.Plan.MY;
    const vertices = new Float32Array(lines.length * 5);
    let index = 0;

    for (let i = 0; i < lines.length; i++) {
        vertices[index++] = lines[i].X;
        vertices[index++] = lines[i].Y;
        vertices[index++] = lines[i].R;
        vertices[index++] = lines[i].G;
        vertices[index++] = lines[i].B;
    }

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

    while (true) {
        await new Promise(r => setTimeout(r, 1));

        if (s > 10) {
            fx = -1;
        } 
        if (s < 0) {
            fx = 1;
        }

        s += 0.01 * fx;

        device.queue.writeBuffer(uniformBuffer, 0, new Float32Array([s, s, 0, 0]));

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
        passEncoder.draw(lines.length);

        passEncoder.end();
        device.queue.submit([commandEncoder.finish()]);
    }
}

window.onload = init;
