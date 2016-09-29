# probe

probe is a simple tool to list containers once docker ps hang due to
a container stuck on staring or stopping.

probe need to root access to root of docker runtime.
The dafault root is `/var/run/docker` and use `-r` to 
specify a custom one. probe display the running container
by default and use `-a` to display all the containers.

Usage:
````
Usage of probe:
  -a    Display all containers
  -r string
        Root of the Docker runtime (default "/var/lib/docker")
````

Example:
````
sudo probe
ID              Image           Created         cmd             Status                  Name
a55a5112fa7e    2b8fd9751c4c    3 weeks         sh              Up 3 weeks              /sharp_perlman
````
