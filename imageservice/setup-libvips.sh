  sudo yum update -y \
  && sudo yum groupinstall -y "Development Tools" \
  && sudo yum install -y wget tar

  sudo yum install -y \
  libpng-devel \
  glib2-devel \
  libjpeg-devel \
  expat-devel \
  zlib-devel 

  sudo yum install -y gcc64

  VIPS_VERSION=8.8.0
  VIPS_URL=https://github.com/libvips/libvips/releases/download
  
  sudo ln -s /usr/lib/gcc/x86_64-amazon-linux/6.4.1/libgomp.spec /usr/lib64/libgomp.spec
  sudo ln -s /usr/lib/gcc/x86_64-amazon-linux/6.4.1/libgomp.a /usr/lib64/libgomp.a
  sudo ln -s /usr/lib64/libgomp.so.1.0.0 /usr/lib64/libgomp.so

  wget ${VIPS_URL}/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz \
  && tar xzf vips-${VIPS_VERSION}.tar.gz \
  && cd vips-${VIPS_VERSION} \
  && ./configure \
  && sudo make \
  && sudo make install