# MCSPyW

The _Minecraft Server Python Wrapper_.

## Gallery
Here are some examples of what can be done using MCSPyW

You can run any existing Minecraft command by appending a ! in front:

    !time set 0
![Setting the time to night](http://i.imgur.com/qZ8sIaH.png)

I've also been experimenting with **plotting functions** in Minecraft.
To plot something, run the following command:

    !plot <function> <length> <material> [material-data]
Where \<function> is a valid python lambda function and \<length> is how long you want the function to run for.

You can plot straight lines:

    !plot (lambda i: (-240 -i, 110 -270)) 100 wool
![Plotting a staight line](http://i.imgur.com/BPJRguN.png)
But that's boring.

Here's a parabola:

    !plot (lambda i: (-240 -i, 110 + (50 - i)**2 / 20, -270)) 100 wool 3
![Plotting a  parabola](http://i.imgur.com/DwtGSBM.png)

And a sin curve:

    !plot (lambda i: (-240 -i, 110 + 10 * sin(i / 5), -270)) 100 wool 1
![Plotting a sin curve](http://i.imgur.com/biixx7F.png)

And a spiral:

    !plot (lambda i: (-240 -i, 110 + 10 * sin(i / 5), -270 + 10 * cos(i / 5))) 100 wool 2
![Plotting a spiral](http://i.imgur.com/o9eclmB.png)
