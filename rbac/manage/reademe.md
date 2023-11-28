# 目标

对策略进行管理，对策略的所有修改都由这里去实现，这里会想办法让其他服务知道策略已经修改并需要重新加载

# 实现

## 方案1

这个实现为管理服务 ，简称 admin
首先每个cacheadapter实例要有一个自己的唯一id，这个id要让admin知道
然后每一个权限认证进行的时候都需要先到redis检测这个值是什么情况，如果需要修改就重新进行load，如果不需要进行修改就不管
方法1就是消息队列
方法2就是 使用特定的前缀将这个唯一id放入redis中去，然后admin需要用到的时候，就去操作这个特定前缀的就可以了

admin持有一个单独的cacheadapter，进行数据库修改之后，用这个唯一id去更改redis里边固定前缀键的值，

值不为load就要重新加载，然后设置为load

# admin服务能做那些修改

无非就是禁止某个用户访问某个资源
让某个用户能访问某个资源