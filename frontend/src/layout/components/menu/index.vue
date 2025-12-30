<template>
  <el-menu
          class="menus"
          :collapse="isCollapse"
          :default-active="activeIndex"
          background-color="transparent"
          text-color="rgba(255, 255, 255, 0.8)"
          active-text-color="#a78bfa"
          router
  >
    <menu-item v-for="(route,index) in routes" :key="route.name" :route="route"></menu-item>
  </el-menu>
</template>

<script>
  import { defineComponent, ref, onMounted, watch, computed } from 'vue'
  import { useRouteStore } from '@/store/router'
  import MenuItem from '@/layout/components/menu/item.vue'
  import { useRoute } from 'vue-router'
  import { useAppStore } from '@/store/app'

  export default defineComponent({
    name: 'Menu',
    created () {
    },
    components: { MenuItem },
    setup () {
      const routes = ref([])
      const route = useRoute()
      const app = useAppStore()
      const isCollapse = computed(() => app.setting.sideIsCollapse)
      const activeIndex = computed(() => route.name)

      routes.value = useRouteStore().routes
      return {
        routes,
        activeIndex,
        isCollapse,
      }
    },

  })
</script>

<style lang="scss" scoped>
  .menus {
    min-height: 100vh;
    border-right: none;
    background: transparent;
    
    &:not(.el-menu--collapse) {
      width: var(--sideBarWidth);
    }

    :deep(.el-menu-item) {
      transition: all var(--transition-base);
      border-left: 3px solid transparent;
      margin: var(--spacing-xs) var(--spacing-sm);
      border-radius: var(--border-radius-sm);
      
      &:hover {
        background: rgba(255, 255, 255, 0.1) !important;
        transform: translateX(4px);
      }
      
      &.is-active {
        background: linear-gradient(90deg, rgba(167, 139, 250, 0.2) 0%, rgba(139, 92, 246, 0.1) 100%) !important;
        border-left-color: #a78bfa;
        font-weight: 500;
      }
    }

    :deep(.el-sub-menu__title) {
      transition: all var(--transition-base);
      border-left: 3px solid transparent;
      margin: var(--spacing-xs) var(--spacing-sm);
      border-radius: var(--border-radius-sm);
      
      &:hover {
        background: rgba(255, 255, 255, 0.1) !important;
        transform: translateX(4px);
      }
    }

    :deep(.el-icon) {
      transition: transform var(--transition-base);
      
      &:hover {
        transform: scale(1.1);
      }
    }
  }
</style>
<style>
</style>
