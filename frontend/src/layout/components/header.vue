<template>
  <el-icon class="ex-icon" @click="expandOrFoldSlider">
    <el-icon-expand v-if="setting.sideIsCollapse"></el-icon-expand>
    <el-icon-fold v-else></el-icon-fold>
  </el-icon>
  <div class="header-logo">
    <img :src="setting.logo" alt="" class="logo">
    <div class="title gradient-text">{{setting.title}}</div>
  </div>
  <Setting></Setting>
</template>

<script>
  import { defineComponent, computed } from 'vue'
  import HeaderMenu from '@/layout/components/menu/index.vue'
  import Setting from '@/layout/components/setting/index.vue'
  import { useAppStore } from '@/store/app'
  import GTags from '@/layout/components/tags/index.vue'

  export default defineComponent({
    name: 'LayerHeader',
    created () {
    },
    components: { HeaderMenu, Setting, GTags },
    watch: {},
    setup (props) {
      const appStore = useAppStore()
      const setting = computed(() => appStore.setting)
      const expandOrFoldSlider = () => {
        appStore.sideCollapse()
      }
      return {
        setting,
        expandOrFoldSlider,
      }
    },

  })
</script>

<style scoped lang="scss">
  .ex-icon {
    height: 100%;
    display: flex;
    align-items: center;
    margin-right: var(--spacing-md);
    font-size: 20px;
    cursor: pointer;
    padding: var(--spacing-sm);
    border-radius: var(--border-radius-sm);
    transition: all var(--transition-base);
    color: var(--text-primary);

    &:hover {
      background: var(--gradient-subtle);
      transform: scale(1.1);
    }

    &:active {
      transform: scale(0.95);
    }
  }

  .header-logo {
    display: flex;
    height: 100%;
    align-items: center;
    gap: var(--spacing-md);

    .title {
      display: block;
      font-weight: 600;
      font-size: 1.25rem;
      letter-spacing: -0.025em;
    }

    .logo {
      display: block;
      width: 36px;
      height: 36px;
      border-radius: var(--border-radius-sm);
      transition: all var(--transition-base);
      animation: glow 3s ease-in-out infinite;

      &:hover {
        transform: rotate(360deg) scale(1.1);
      }
    }
  }

</style>
<style lang="scss">

</style>
