#
# Cookbook Name:: dd-agent-sles-workaround
# Recipe:: default
#
# Copyright (C) 2020 Datadog
#
# All rights reserved - Do Not Redistribute
#

if node['platform_family'] == 'suse'
  # Stop the Windows Azure Agent, for some reason it's changing the hostname
  # every 30s on SLES 11 and 12, which leads to a network interface reset,
  # making it likely for tests to fail if a network call happens during the reset.
  service 'waagent' do
    service_name "waagent"
    action :stop
  end

  # Wait a bit to let the Azure Agent go down
  execute 'Wait a bit to let Azure Agent stop' do
    command "sleep 5"
  end
end
