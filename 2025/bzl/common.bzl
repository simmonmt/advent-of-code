def combine_samples(name, srcs, out):
    native.genrule(
        name = name,
        srcs = srcs,
        outs = [out],
        cmd = "(for file in $(SRCS) ; do " +
              "  echo $$file ; cat $$file |wc -l ; cat $$file ; " +
              "done) >$@",
    )
