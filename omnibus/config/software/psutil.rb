name "psutil"
default_version "5.6.7"

version "5.6.7" do
  source sha256: "ebbed18bf912fa4981d05c5c1a0cacec2aa8c594e1608f0cf5cc7a3d4f63d4d4"
end

dependency "python"
dependency "pip"

source url: "https://github.com/giampaolo/psutil/archive/release-#{version}.tar.gz"

relative_path "psutil-release-#{version}"

env = with_embedded_path
env = with_standard_compiler_flags(env)

if linux?
  env = with_glibc_version(env)
  env['CFLAGS'] = "-D_DISABLE_PRLIMIT #{env['CFLAGS']}"
end

build do
  ship_license "https://raw.githubusercontent.com/giampaolo/psutil/master/LICENSE"

  patch source: "psutil-5.6.7-hackadog.patch", env: env

  pip "install --install-option=\"--install-scripts=#{windows_safe_path(install_dir)}/bin\" .", :env => env
end
