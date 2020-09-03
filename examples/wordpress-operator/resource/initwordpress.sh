#!/bin/sh

curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar
chmod +x wp-cli.phar
mv wp-cli.phar /usr/local/bin/wp

wp --info


wp core is-installed --allow-root
status=$?
echo $status
path="wordpress-$5"
url=$6
if [ "$status" -eq "1" ]; then
    wp core install --url=$url --title=$1 --admin_user=$2 --admin_password=$3 --admin_email=$4 --allow-root

    
    cat <<PHP >> /var/www/html/wp-config.php

\$host = '$url';
define('WP_HOME', \$host);
define('WP_SITEURL', \$host);
PHP

fi

