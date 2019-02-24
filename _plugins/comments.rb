require 'digest'
require 'json'
require 'open3'

# Adapted from https://github.com/aleung/jekyll-post-revision

module Jekyll
  module Comments

    class Generator < Jekyll::Generator
      def generate(site)
        return if ARGV.include?("--no-comments")

        sha2 = Digest::SHA2.new

        site.posts.each do |post|

          # Create a post ID based on post date
          sha2.reset
          sha2.update post.data['date'].strftime('%s')
          post_id = sha2.hexdigest

          # Download comments for post
          # TODO: Place URL's in site_config
          # TODO: Put plugin version number in user agent string
          # TODO: Maybe using backticks is enough ( raw = `curl [...]` )
          raw = Executor.sh('curl', '-Ssf',
                            '-H', 'Accept: application/json',
                            '-H', 'User-Agent: jekyll-express-comments',
                            'https://blog.stigok.com/comments/' + post_id)

          comments = JSON.parse(raw.lines.join('\n'))

          post.data['post_id'] = post_id
          post.data['comments'] = comments

          puts "%d comments for post id %s" % [comments.count, post_id]
        end
      end
    end # Generator

    module Executor
      def self.sh(*args)
        Open3.popen2e(*args) do |stdin, stdout_stderr, wait_thr|
          exit_status = wait_thr.value # wait for it...
          output = stdout_stderr.read
          output ? output.strip : nil
        end
      end
    end # Executor

  end # Comments
end # Jekyll

