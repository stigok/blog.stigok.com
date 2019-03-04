require 'digest'

# Add a subject ID to all posts
Jekyll::Hooks.register :posts, :pre_render do |post|
  post_hash = Digest::SHA256.hexdigest post.data['date'].strftime('%s')
  post.data['comments_subject_id'] = post_hash
end

module Jekyll
  class PageWithoutAFile < Jekyll::Page
    def read_yaml(*)
      @data ||= {}
    end
  end

  # Create a list of subject IDs for the comment web server to validate against
  class CommentsSubjectListGenerator < Generator
    safe true
    priority :lowest

    def generate(site)
      filename = 'comments_subject_ids.json'

      mapping = {}
      site.posts.each do |post|
        hash = Digest::SHA256.hexdigest post.data['date'].strftime('%s')
        mapping[hash] = post.url
      end

      page = PageWithoutAFile.new(site, __dir__, '', filename).tap do |file|
        file.content = JSON.generate(mapping)
        file.output
      end

      site.pages << page
    end
  end
end
