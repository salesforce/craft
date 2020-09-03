# Tutorial : Wordpress Operator
Unlike most tutorials who start with some really contrived setup, or some toy application that gets the basics across, this tutorial will take you through the full extent of the CRAFT application and how it is useful. We start off simple and in the end, build something pretty full-featured and meaningful, namely, an operator for the Wordpress application.

The job of the Wordpress operator is to host the Wordpress application and perform operations given by the user on the cluster. It reconciles regularly, checking for updates in the resource and therefore, can be termed as level-triggered.

We will see how the controller.json and resource.json required to run the wordpress operator have been developed. Then, we'll see how to use CRAFT to create the operator and deploy it onto the cluster.

The config files required to create the operator are already present in `example/wordpress-operator`

Let's go ahead and see how we have created our files for the wordpress application to understand and generalise this process for any application. First, we start with the controller.json file in the next section.

