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
      # TODO: Figure out how to make this generator run after the hook on top, so that we don't
      #       have to calculate the hash again.
      subject_ids = site.posts.map { |post| Digest::SHA256.hexdigest post.data['date'].strftime('%s') }

      page = PageWithoutAFile.new(site, __dir__, '', filename).tap do |file|
        file.content = JSON.generate(subject_ids)
        file.output
      end

      site.pages << page
    end
  end
end
