/**
 * Created by zzmhot on 2017/3/23.
 *
 * 路由Map
 *
 * @author: zzmhot
 * @github: https://github.com/zzmhot
 * @email: zzmhot@163.com
 * @Date: 2017/3/23 18:30
 * @Copyright(©) 2017 by zzmhot.
 *
 */

import Vue from 'vue'
import VueRouter from 'vue-router'
import store from 'store'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

Vue.use(VueRouter)

//使用AMD方式加载
// component: resolve => require(['pages/home'], resolve),
const routes = [{
        path: '/home',
        name: 'home',
        components: {
            default: require('pages/home'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "主页",
            auth: true
        }
    }, {
        path: '/conf/list',
        name: 'confList',
        components: {
            default: require('pages/conf/base'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "项目配置",
            auth: true
        }
    }, {
        path: '/conf/update/:id',
        name: 'confUpdate',
        components: {
            default: require('pages/conf/save'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "修改配置",
            auth: true
        }
    }

    , {
        path: '/conf/detection/:id',
        name: 'confDetection',
        components: {
            default: require('pages/conf/detection'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "检测项目",
            auth: true
        }
    }

    , {
        path: '/conf/add',
        name: 'confAdd',
        components: {
            default: require('pages/conf/save'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "增加项目配置",
            auth: true
        }
    }, {
        path: '/task/list',
        name: 'taskList',
        components: {
            default: require('pages/task/base'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "全部上线单",
            auth: true
        }
    }, {
        path: '/task/searchlist',
        name: 'searchtaskList',
        components: {
            default: require('pages/task/searchbase'),
            menuView: require('components/leftSlideTologin')

        }
    }, {
        path: '/task/create',
        name: 'taskCreate',
        components: {
            default: require('pages/task/create'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "创建上线单",
            auth: true
        }
    }, {
        path: '/task/release/:id',
        name: 'taskRelease',
        components: {
            default: require('pages/task/release'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "部署上线",
            auth: true
        }
    }, {
        path: '/task/searchrelease/:id',
        name: 'searchtaskRelease',
        components: {
            default: require('pages/task/searchrelease'),
            menuView: require('components/leftSlideTologin')
        },
        meta: {
            title: "部署上线",
            auth: true
        }
    }, {
        path: '/task/mylist',
        name: 'taskMyList',
        components: {
            default: require('pages/task/mylist'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "我的上线单",
            auth: true
        }
    }, {
        path: '/task/git',
        name: 'taskGit',
        components: {
            default: require('pages/task/git'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "创建上线单",
            auth: true
        }
    }, {
        path: '/task/file',
        name: 'taskFile',
        components: {
            default: require('pages/task/file'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "创建上线单",
            auth: true
        }
    }, {
    path: '/task/jenkins',
    name: 'taskJenkins',
    components: {
      default: require('pages/task/jenkins'),
      menuView: require('components/leftSlide')
    },
    meta: {
      title: "创建上线单",
      auth: true
    }
  }, {
        path: '/p2p/check',
        name: 'p2pCheck',
        components: {
            default: require('pages/p2p/check'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "agent状态查询",
            auth: true
        }
    }, {
        path: '/other/noauto',
        name: 'noauto',
        components: {
            default: require('pages/other/noauto'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "预发布统计",
            auth: true
        }
    }, {
        path: '/other/flush',
        name: 'otherFlush',
        components: {
            default: require('pages/task/flush'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "刷新版本号",
            auth: true
        }
    }, {
        path: '/other/gitpull',
        name: 'otherGitpull',
        components: {
            default: require('pages/task/gitpull'),
            menuView: require('components/leftSlide')
        },
        meta: {
            title: "预发布git版本查看",
            auth: true
        }
    }, {
        path: '/user/login',
        name: 'login',
        components: {
            fullView: require('pages/user/login')
        }
    }, {
        path: '/user/register',
        name: 'register',
        components: {
            menuView: require('components/leftSlide'),
            default: require('pages/user/register')

        }
    }, {
        path: '/user/changepasswd',
        name: 'changepasswd',
        components: {
            menuView: require('components/leftSlide'),
            default: require('pages/user/changepasswd')

        },
        meta: {
            title: "修改密码",
            auth: true
        }
    }, {
    path: '/user/list',
    name: 'userList',
    components: {
      menuView: require('components/leftSlide'),
      default: require('pages/user/list')

    },
    meta: {
      title: "用户列表",
      auth: true
    }
  }, {
        path: '',
        redirect: '/home'
    }, {
        path: '*',
        name: 'notPage',
        components: {
            fullView: require('pages/error/404')
        }
    }
]

const router = new VueRouter({
    routes,
    mode: 'hash', //default: hash ,history
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { x: 0, y: 0 }
        }
    }
})

//全局路由配置
//路由开始之前的操作
router.beforeEach((to, from, next) => {
    NProgress.start()
    let toName = to.name
        // let fromName = from.name
    let is_login = store.state.user_info.login


    if (!is_login && toName === 'searchtaskList') {
        next();
    } else if (!is_login && toName === 'searchtaskRelease') {
        next({});
    } else if (!is_login && toName !== 'login') {
        next({
            name: 'login'
        });
    } else {
        if (is_login && toName === 'login') {
            next({
                path: '/task/list'
            });
        } else {
            next();
        }
    }
})

//路由完成之后的操作
router.afterEach(route => {
    NProgress.done()
})

export default router
