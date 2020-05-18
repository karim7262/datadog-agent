require 'spec_helper'

describe 'dd-agent-installation-script' do
  include_examples 'Agent'

  context 'when testing DD_SITE' do
    let(:config) do
      YAML.load_file('/etc/datadog-agent/datadog.yaml')
    end

    it 'uses DD_SITE to set the site' do
      expect(config['site']).to eq 'datadoghq.eu'
    end
  end
end
