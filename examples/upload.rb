require 'rest_client'

r = RestClient.post("http://localhost:8022/v1/files?token=hello", file: File.new('upload.rb'))
p r
