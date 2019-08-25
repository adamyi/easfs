'use strict';

/* eslint browser:true */

const body = document.querySelector('body');
const isDocPage = document
  .querySelector('body.devsite-doc-page') ? true : false;

function highlightActiveNavElement() {
  var elems = document.querySelectorAll('.devsite-section-nav li.devsite-nav-active');
  for (var i = 0; i < elems.length; i++) {
    expandPathAndHighlight(elems[i]);
  }
}

function expandPathAndHighlight(elem) {
  // Walks up the tree from the current element and expands all tree nodes
  var parent = elem.parentElement;
  var parentIsCollapsed = parent.classList.contains('devsite-nav-section-collapsed');
  if (parent.localName === 'ul' && parentIsCollapsed) {
    parent.classList.toggle('devsite-nav-section-collapsed');
    parent.classList.toggle('devsite-nav-section-expanded');
    // Checks if the grandparent is an expandable element
    var grandParent = parent.parentElement;
    var grandParentIsExpandable = grandParent.classList.contains('devsite-nav-item-section-expandable');
    if (grandParent.localName === 'li' && grandParentIsExpandable) {
      var anchor = grandParent.querySelector('a.devsite-nav-toggle');
      anchor.classList.toggle('devsite-nav-toggle-expanded');
      anchor.classList.toggle('devsite-nav-toggle-collapsed');
      expandPathAndHighlight(grandParent);
    }
  }
}

function getCookieValue(name, defaultValue) {
  const value = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)');
  return value ? value.pop() : defaultValue;
}

function initYouTubeVideos() {
  var videoElements = body
    .querySelectorAll('iframe.devsite-embedded-youtube-video');
  videoElements.forEach(function(elem) {
    const videoID = elem.getAttribute('data-video-id');
    if (videoID) {
      let videoURL = 'https://www.youtube.com/embed/' + videoID;
      videoURL += '?autohide=1&amp;showinfo=0&amp;enablejsapi=1';
      elem.src = videoURL;
    }
  });
}

function init() {
  initYouTubeVideos();
  highlightActiveNavElement();
}

init();
