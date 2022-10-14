docker volume prune
docker system prune

function update_images {
    echo "==> Updating Images"
    IMAGES=$(docker images | tail -n +2 | awk '{ print $1 ":" $2 }' | sort)
    for I in $IMAGES; do
        [ ! "$I" == "<none>:<none>" ] && docker pull $I || echo "Skipping untagged image"
    done
}

function rm_dangling {
    echo "==> Deleting Dangling Images"
    IMAGES=$(docker images -q --filter dangling=true)
    test "$IMAGES" && (echo $IMAGES | xargs -n 1 docker rmi) || true
}

function rm_untagged {
    echo "==> Deleting Untagged Images"
    IMAGES=$(docker images | grep '<none>' | awk '{ print $3 }')
    test "$IMAGES" && (echo $IMAGES | xargs -n 1 docker rmi) || true
}

function rm_exited {
    echo "==> Deleting Stopped Containers"
    CONTAINERS=$(docker ps -q --filter status=exited --filter status=created --filter status=dead)
    test "$CONTAINERS" && (echo $CONTAINERS | xargs -n 1 docker rm -f ) || true
}

function rm_volumes {
    echo "==> Deleting Danging Volumes"
    VOLUMES=$(docker volume ls -qf dangling=true)
    test "$VOLUMES" && (echo $VOLUMES | xargs -n 1 docker volume rm) || true
}

# shellcheck disable=SC2120
function rm_versioned {
    echo "==> Deleting Versioned Images"
    IMAGE="$1"
    KEEP=${2:-4}

    [ "$KEEP" -eq "$KEEP" ] || (echo "Error KEEP isnt a number" && exit 128)
    [ "$IMAGE" == "" ] && echo "Usage: $0 <image> <keep>" && exit 1

    # skip first line, show full image name, sort, then grab $KEEP number
    IMAGES=$(docker images "$IMAGE" | tail -n+2 | awk '{ print $1 ":" $2 }' | sort --version-sort --reverse | tail -n+$((1 + $KEEP)))

    for I in $IMAGES; do
       docker rmi $I 1> /dev/null
    done
}

function restart_docker {
  systemctl restart docker.service
}

function check_docker {
  docker build -<<EOF
  FROM busybox
  RUN echo "hello world"
EOF
}

function clean_builder_cache {
  docker builder prune
}

update_images
rm_dangling
rm_untagged
rm_exited
rm_volumes
rm_versioned
clean_builder_cache