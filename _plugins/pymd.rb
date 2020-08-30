require 'open3'

module Jekyll
  module Pymd

    class Generator < Jekyll::Generator
      def generate(site)
        return if ARGV.include?("--no-pymd")

        @site = site
        site.posts.docs.each do |post|
          # Only run if pymd is explicitly specified
          next unless post.data.has_key?('processors')
          next unless post.data['processors'].include?('pymd')

          puts "pymd processing: \"#{post.data['title']}\""
          stdout, res = Open3.capture2("python #{__dir__}/pymd.py -", :stdin_data=>post.content)

          if res != 0 then
            #raise "pymd failed: #{res}"
            Process.exit 1
          end

          post.content = stdout
          post
        end
      end
    end

  end
end
