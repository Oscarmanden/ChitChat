
Exercise - Physical time with gRPC

    Create a a gRPC service endpoint created in Golang, that returns the current time (from time.Now())
    Each member of your group should deploy the time service node implementation on their computer.
    Expose your service to other group members - assuming the group is on the same network. 
        Format your port for TCP with ":<portnumber>".
        If you have trouble with the network you can try running multiple services on the same device device - however you might not see the differences in time between devices if you do so.
        Note: Do you need to change the firewall settings on your laptop? What network command can you use to check, if you can access the port on your friends computer?
    Create a TimeServiceClient, that connects to at least 2 different service nodes (possibly on the different team members' computers) that read the time and output it to a log file. Can you observe an offset? Can you observe clock skew?
    Enhance the client, so that it can measure the round-trip time for each service call.
    Optional: Discuss Christian's Synchronization Algorithm in the context of your time service client. Can you implement it?

This exercise is very relevant for your next mandatory hand-in ✨🦄


Optional Exercise - Use containers for Golang and Protobuf development

Note! This exercise is optional.

If you are up for an extra challenge you could try to work with containers. You are free to use Docker, Podman or something else entirely. If you have not already you are very likely to run into containers at some point in your software career.

You do not need to build a Distributed System to see the benefits of using containers. However, it is all up to you and you are not forced to use containers at all.

One of the benefits of containers is to keep your system clean, and avoid having to install many applications natively and keep track of which versions you have installed. With containers you can get the dependencies you need for the specific version of your software without those version having to match your other projects. There can be several reasons why keeping a specific version of a dependency on your device natively can be troublesome.

Start by containerizing your Go development, so you don't use the compiler natively installed on your computer anymore:

    Go through the contents of the Dockerfile created by Chris Crone, which lets you containerize your Go development. You can check out blog posts 1-3, where he explains the contents and usage: https://www.docker.com/blog/tag/go-env-series/.
    Let us see if you can use Docker to compile your proto files from last weeks exercise set? You can try to use an already existing Docker image, or add a multistage build step to the Dockerfile you created in step 1.

There is a reference implementation here.