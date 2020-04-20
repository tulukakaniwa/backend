FROM centos:centos7 AS base-dependencies

# install gcc 7
RUN yum -y install centos-release-scl && \
    yum -y install devtoolset-7 devtoolset-7-libatomic-devel
ENV CC=/opt/rh/devtoolset-7/root/usr/bin/gcc \
    CPP=/opt/rh/devtoolset-7/root/usr/bin/cpp \
    CXX=/opt/rh/devtoolset-7/root/usr/bin/g++ \
    PATH=/opt/rh/devtoolset-7/root/usr/bin:$PATH \
    LD_LIBRARY_PATH=/opt/rh/devtoolset-7/root/usr/lib64:/opt/rh/devtoolset-7/root/usr/lib:/opt/rh/devtoolset-7/root/usr/lib64/dyninst:/opt/rh/devtoolset-7/root/usr/lib/dyninst:/opt/rh/devtoolset-7/root/usr/lib64:/opt/rh/devtoolset-7/root/usr/lib:$LD_LIBRARY_PATH

# python 3.6
RUN yum -y install rh-python36
ENV PATH=/opt/rh/rh-python36/root/usr/bin:$PATH


# install other yosys dependencies
RUN yum install -y flex tcl tcl-devel libffi-devel git graphviz readline-devel glibc-static wget autoconf zlib-devel && \
    wget https://ftp.gnu.org/gnu/bison/bison-3.0.1.tar.gz && \
    tar -xvzf bison-3.0.1.tar.gz && \
    cd bison-3.0.1 && \
    ./configure && \
    make -j$(nproc) && \
    make install

FROM base-dependencies AS builder

COPY build/yosys /yosys
WORKDIR /yosys
RUN make config-gcc-static-tcl-dynamic
RUN make -j$(nproc)
RUN make install

# TODO: copy openroad application

FROM centos:centos7 as runner

RUN yum update -y && \
    tum install -y tcl tcl-devel libffi-devel readline-devel

COPY --from=builder /usr/local/bin/yosys /usr/local/bin/yosys
COPY --from=builder /usr/local/bin/yosys-abc /usr/local/bin/yosys-abc

# TODO: copy required binaries