# Usage:
# docker run --rm -it -v $(pwd):/workdir xavisb/dog -i task-name

FROM alpine:3.5

# Requires the binary to be already built
COPY dist/linux_amd64/dog /usr/local/bin/dog

# Copy or mount your project and Dogfiles at /workdir
WORKDIR /workdir

ENTRYPOINT ["dog"]
