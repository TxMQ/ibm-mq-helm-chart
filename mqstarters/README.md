Starter projects are based on the idea of application boundary.

Application boundary is defined by a set of channels and objects that are needed for the application.

To make this possible an application is allocated application object tree.

For example, if application is called app1 then object tree is rooted at appl1.

For each application 3 roles are defined: application role, application admin role, and application admin reader role.
Application role grants all required authorities needed to accomplish it's work.
Application admin role grants admin authorities on objects in the appliction tree.
Application admin reader role grants display authorities on objects in the application tree.

This model can be adjusted as needed to extend it or make it more granular.

The template for application boundary.

create channels, use channel names in application tree.
set channel rules
grant connect authority to queue manager

create queues, use queue names in the application tree
create topics, use topic names in the application tree

grant authorities for application role
grant authorities for application admin role
grant authorities for application admin reader role

All authorities are scoped to application tree.

Starter for mq explorer follows this pattern.
