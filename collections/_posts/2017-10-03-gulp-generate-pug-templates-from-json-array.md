---
layout: post
title: "gulp generate pug templates from JSON array"
date: 2017-10-03 16:13:51 +0200
categories: gulp pug
redirect_from:
  - /post/gulp-generate-pug-templates-from-json-array
---

I have a JSON array containing information for different sub pages for a website. I'm already using gulp for the build process, and I want to compile a template pug file into multiple pages based on this JSON file.

    
    const plugins = require('gulp-load-plugins')()
    const gulp = require('gulp')
    const es = require('event-stream')
    const projects = require('./src/projects.json')

    gulp.task('projects', () => {
      // Walk through all projects in JSON
      let streams = projects.map(project => {
        return gulp.src('src/templates/project.pug')
          // Set project as a locals object for template generation
          .pipe(plugins.data(file => ({project: project})))
          .pipe(plugins.pug())
          // Name the file after the project
          .pipe(plugins.rename(path.join(project.name + '.html')))
      })
      // Merge all streams and output to files
      es.merge(streams)
        .pipe(gulp.dest(path.join(paths.dist, 'projects')))
    })