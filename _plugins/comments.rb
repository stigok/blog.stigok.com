require 'digest'

# Add a subject ID to all posts
Jekyll::Hooks.register :posts, :pre_render do |post|
  post_hash = Digest::SHA256.hexdigest post.data['date'].strftime('%s')
  post.data['comments_subject_id'] = post_hash
end

