#!/bin/bash

TEMP_FILE='tempFile.txt'

getResource() {
  echo $1
  curl --compressed --fail --progress-bar $1 > $TEMP_FILE
  if [ $? -eq 0 ]; then
      mv $TEMP_FILE $2
  else
      echo Unable to update $1
      rm $TEMP_FILE
  fi
}

echo "Getting lastest CSS and Script resources..."

getResource https://developers.google.cn/_static/css/devsite-google-blue.css static/styles/devsite-google-blue.css
getResource https://developers.google.cn/_static/css/devsite-orange.css static/styles/devsite-orange.css
getResource https://developers.google.cn/_static/js/script_foot_closure.js static/scripts/script_foot_closure.js
getResource https://developers.google.cn/_static/js/framebox.js static/scripts/framebox.js
getResource https://developers.google.cn/_static/js/jquery-bundle.js static/scripts/jquery-bundle.js
getResource https://developers.google.cn/_static/js/prettify-bundle.js static/scripts/prettify-bundle.js

