#
# Copyright:: Copyright (c) 2013-2014 Chef Software, Inc.
# License:: Apache License, Version 2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

name "jemalloc"

default_version "5.2.0"

source :url => "https://github.com/jemalloc/jemalloc/archive/#{version}.tar.gz",
       :sha256 => "acd70f5879700567e1dd022dd11af49100c16adb84555567b85a1e4166749c8d"

env = {
  "LDFLAGS" => "-L#{install_dir}/embedded/lib -I#{install_dir}/embedded/include",
  "CFLAGS" => "-L#{install_dir}/embedded/lib -I#{install_dir}/embedded/include",
  "LD_RUN_PATH" => "#{install_dir}/embedded/lib",
  # tell the Makefile which is the directory containing config files by setting
  # `conf_dir`, otherwise `make install` will write to `/etc/`
  "conf_dir" =>  "#{install_dir}/embedded/etc"
}

build do
  ship_license "TODO"

  python_configure = ["./configure",
                      "--prefix=#{install_dir}/embedded"]

  command python_configure.join(" "), :env => env
  command "make -j #{workers}", :env => env
  command "make install", :env => env

end
