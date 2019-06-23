# Blox

Blox wraps the Minecraft server executable, and acts as a proxy between the Minecraft client and server. It listens to the processes STDOUT and sniffs the network traffic from the client to extract any useful information we can use to create custom plugins with. Then we either write to the server's STDIN or write directly to the servers network connection.

This way, we can add custom functionality to Minecraft without modifying the code of the Minecraft server.

## Gallery

Here are some examples of what can be done using Blox:

Set the time to day / night.

```
!time day
```

![Setting the time to day](https://imgur.com/esRvpo1.png)

```
!time night
```

![Setting the time to night](https://imgur.com/KSd2Wlx.png)

I've also been experimenting with plotting functions in Minecraft. To do this I have created a custom plotting language which gets evaluated by the wrapper. To plot something, run a command with the following syntax:

```
!plot MATERIAL X_EXPRESSION, Y_EXPRESSION, Z_EXPRESSION | PARAMETER <- RANGE_START..RANGE_END
```

You can have multiple comma separated ranges as you'll see in some of the examples below.

You can plot straight lines:

```
!plot stone 320 + x, 98, 155 | x <- 1..100
```

![Plotting a staight line](https://imgur.com/pbOEr4Z.png)

But that's boring.

Here's a parabola:

```
!plot stone 320 + x, 98 + (x - 50)^2 / 20, 155 | x <- 1..100
```

![Plotting a  parabola](https://imgur.com/KjZbV92.png)

And a spiral:

```
!plot stone 320 + x, 98 + 10 * sin(x / 5), 155 + 10 * cos(x / 5) | x <- 1..100
```

![Plotting a spiral](https://imgur.com/oNhcfQy.png)

How about that cool heart curve your maths teacher always shows you on Valentine's day?

```
!plot stone 320 + i / 20 * (16 * sin(t) ^ 3), 130 + i / 20 * (13 * cos(t) - 5 * cos(2 * t) - 2 * cos( 3 * t) - cos(4 * t)), 155 | t <- 1..1000, i <- 1..60
```

![Plotting a heart](https://imgur.com/hSIV5Br.png)

We can even get 3 dimensional!

Here's a Paraboloid:

```
!plot stone 320 + x, 98 + ( (x - 20)^2 + (y - 20)^2 ) / 20, 155 + y | x <- 1..40, y <- 1..40
```

![Plotting a paraboloid](https://imgur.com/FzPM1UV.png)

And a ripple curve thing:

```
!plot stone 320 + x, 108 + 5 * sin( ((x - 30)/10)^2 + ((y - 30)/10)^2 ), 155 + y | x <- 1..60, y <- 1..60
```

![Plotting a ripple curve](https://imgur.com/xVSYSK2.png)

This one makes me wonder if we can use this tool to generate interesting landscapes from within Minecraft:

```
!plot stone 320 + x, 108 + 10 * sin((x - 30)/5) * cos((y - 30)/5), 155 + y | x <- 1..60, y <- 1..60
```

![Plotting a wavey curve](https://imgur.com/fWdm7Nm.png)
