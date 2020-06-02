def _convert_md_to_html_impl(ctx):
    outs = []
    for in_f in ctx.files.srcs:
        out_f = ctx.actions.declare_file(in_f.basename + ".html")
        outs.append(out_f)
        ctx.actions.run(
            inputs = [ctx.file.template, in_f],
            outputs = [out_f],
            arguments = ["--in", in_f.path, "--out", out_f.path, "--template", ctx.file.template.path],
            executable = ctx.executable._md,
        )
    return DefaultInfo(files = depset(items = outs))

convert_md_to_html = rule(
    implementation = _convert_md_to_html_impl,
    attrs = {
        "srcs": attr.label_list(mandatory = True, allow_files = True),
        "template": attr.label(mandatory = True, allow_single_file = True),
        "_md": attr.label(
            executable = True,
            cfg = "host",
            allow_files = True,
            default = Label("//tools:md"),
        ),
    },
)
