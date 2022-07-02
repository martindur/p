#!/bin/bash

# Helper function for simple navigation of projects
function pcd() {
    if [ -z "$1" ]; then
        project=$(p _getppath)
    else
        project=$(p _getprojectpath $1)
    fi
    cd $project
}
