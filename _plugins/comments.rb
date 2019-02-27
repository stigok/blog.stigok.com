require 'digest'

# Add an ID to all posts based on post date
Jekyll::Hooks.register :posts, :pre_render do |post|
  post.data['post_id'] = Digest::SHA256.hexdigest post.data['date'].strftime('%s')
end

module Jekyll
  module Comments
    class Generator < Jekyll::Generator
      def generate(site)
        return
        site.posts.each do |post|
          id = post.data['post_id']
          comments = site.collections.comments.select{ |com| com[:subject_id] == id }
          post.data['comments'] = comments
          puts "%d comments for post id %s" % [comments.count, id]
        end
      end
    end # Generator
  end # Comments
end # Jekyll
