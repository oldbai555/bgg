<template>
  <div>
    <a-tooltip :trigger="['focus']" placement="topLeft" overlay-class-name="numeric-input">
      <a-input
          type="number"
          :value="val"
          :placeholder="placeholder"
          :max-length="16"
          allow-clear
          @blur="onBlur"
          @change="onChange"
          @pressEnter="onPressEnter"
          :disabled="disabled"
      />
    </a-tooltip>
  </div>
</template>

<script setup lang="ts">
import {defineEmits, defineProps, withDefaults} from "vue";

interface Props {
  val?: number,
  placeholder: string,
  disabled?: boolean,
}

const props = withDefaults(defineProps<Props>(), {
  numb: undefined,
  placeholder: "请输入",
  disabled: false,
})

const emit = defineEmits(['blur', 'change', 'pressEnter', 'update:val'])

const onBlur = () => {
  emit('blur');
};

const onChange = (e: InputEvent) => {
  emit('update:val', Number((e.target as HTMLInputElement).value) || 0);
  emit('change', e);
}

const onPressEnter = (e: InputEvent) => {
  emit('pressEnter', e);
  onChange(e);
  onBlur();
}

</script>

<style scoped>

</style>