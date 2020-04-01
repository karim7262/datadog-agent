def list_files()
  require 'find'
  exclude = [
    'c:/windows/temp',
    'c:/windows/prefetch',
    'c:/windows/installer',
    'c:/windows/winsxs',
    'c:/windows/winsxs',
    'c:/windows/servicing/'
  ].each { |e| e.downcase! }
  return Find.find('c:/windows/').reject { |f| f.downcase.start_with?(exclude) }
end
