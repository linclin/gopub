<template>
    <section class="body-warp">
        <transition name="fade" mode="out-in">
            <router-view name="fullView" style="z-index: 200"></router-view>
        </transition>
        <transition name="fade" mode="out-in">
            <router-view name="menuView" style="z-index: 150"></router-view>
        </transition>
        <main-content style="z-index: 100">
            <transition name="fade" mode="out-in">
                <router-view></router-view>
            </transition>
        </main-content>
    </section>
</template>

<script type="text/javascript">
    import {mainContent} from 'components'

    export default {
        name: 'app',
        components: {
            mainContent
        }
    }
</script>

<style lang="scss" type="text/css" rel="stylesheet/scss">
    @import '~assets/scss/main';

    .fade-enter-active,
    .fade-leave-active {
        transition: all .2s ease;
    }

    .fade-enter,
    .fade-leave-active {
        opacity: 0;
    }
</style>
