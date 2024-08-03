# What is progress.go?

- It's a golang version built with Gin of [progressd.io](https://github.com/fehmicansaglam/progressd.io).
- Support non-ascii characters such as Chinese.

## Free Service
- Progress Bar https://progress.microsvc.xyz/bar/:progress
- Progress Pie https://progress.microsvc.xyz/pie/:progress

## Parameters

- Bar
    |  name  |  place  |  type  |  default  |  comment  |  required  |
    | :----  | :----  | :----  | :----  | :----  | :----  |
    | progress  | path | int | Null | value of progress |  ✅ |
    | title  | query | string | "" | left side title string | ❌ |
    | color  | query | string | 428bca | background color of title string | ❌ |
    | scale  | query | int | 100 | maximum value of progress | ❌ |
    | suffix  | query | string | % | suffix of progress value | ❌ |
    | prefix  | query | string | "" | prefix of progress value | ❌ |
    | width  | query | int | 90 | total progress width | ❌ |
    | height  | query | int | 20 | total progress height | ❌ |
    | fontsize  | query | int | 11 | font size of progress text | ❌ |

- Pie
    |  name  |  place  |  type  |  default  |  comment  |  required  |
    | :----  | :----  | :----  | :----  | :----  | :----  |
    | progress  | path | int | Null | value of progress |  ✅ |
    | size  | query | int | 17 | diameter of pie | ❌ |
    | scale  | query | int | 100 | maximum value of progress | ❌ |
    | suffix  | query | string | % | suffix of progress value | ❌ |
    | prefix  | query | string | "" | prefix of progress value | ❌ |
    | fontsize  | query | int | 11 | font size of progress text | ❌ |

## Examples

#### Bar

https://progress.microsvc.xyz/bar/28
![Progress](https://progress.microsvc.xyz/bar/28)

https://progress.microsvc.xyz/bar/28?title=progress
![Progress](https://progress.microsvc.xyz/bar/28?title=progress)   

https://progress.microsvc.xyz/bar/58
![Progress](https://progress.microsvc.xyz/bar/58)   

https://progress.microsvc.xyz/bar/59?title=completed&color=af0000
![Progress](https://progress.microsvc.xyz/bar/58?title=completed&color=af0000)  

https://progress.microsvc.xyz/bar/91?width=300
![Progress](https://progress.microsvc.xyz/bar/91?width=300)  

https://progress.microsvc.xyz/bar/91?title=done
![Progress](https://progress.microsvc.xyz/bar/91?title=done)   

https://progress.microsvc.xyz/bar/7?scale=10&title=mark&suffix=X
![Progress](https://progress.microsvc.xyz/bar/7?scale=10&title=mark&suffix=X)

https://progress.microsvc.xyz/bar/1500?width=500&title=abc&scale=2000&suffix=/$2000&prefix=$
![Progress](https://progress.microsvc.xyz/bar/1500?width=500&title=abc&scale=2000&suffix=/$2000&prefix=$)

#### Pie

https://progress.microsvc.xyz/pie/28
![Progress](https://progress.microsvc.xyz/pie/28)

https://progress.microsvc.xyz/pie/58
![Progress](https://progress.microsvc.xyz/pie/58)    

https://progress.microsvc.xyz/pie/91?size=40&fontsize=40
![Progress](https://progress.microsvc.xyz/pie/91?size=40&fontsize=40)

https://progress.microsvc.xyz/pie/7?scale=10&suffix=X
![Progress](https://progress.microsvc.xyz/pie/7?scale=10&suffix=X)

https://progress.microsvc.xyz/pie/1500?scale=2000&suffix=/$2000&prefix=$
![Progress](https://progress.microsvc.xyz/pie/1500?scale=2000&suffix=/$2000&prefix=$)

---

Heavily inspired by the works of https://github.com/fehmicansaglam/progressd.io and https://github.com/fredericojordan/progress-bar
