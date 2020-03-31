#
# Cookbook Name:: dd-agent-system-files-check
# Recipe:: default
#
# Copyright (C) 2020 Datadog
#
# All rights reserved - Do Not Redistribute
#

if node['platform_family'] != 'windows'
    puts "dd-agent-system-files-check: Not implemented on non-windows"
else
    ruby_block "list-after-files" do
        block do
            File.open("C:/after-files.txt", "w") do |out|
                Dir.glob("C:/windows/**/*").each { |file| out.puts(file) }
            end
        end
        action :run
    end
end
