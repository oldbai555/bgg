<template>
  <div>
    <a-upload
        v-model:file-list="fileList"
        :accept="accept"
        name="file"
        :disabled="disabled"
        :max-count="1"
        :custom-request="customRequest">
      <a-button>
        <upload-outlined></upload-outlined>
        点击上传
      </a-button>
    </a-upload>
  </div>
</template>

<script setup lang="ts">
import {UploadOutlined} from '@ant-design/icons-vue';
import type {UploadProps} from 'ant-design-vue';
import {defineEmits, defineProps, ref, withDefaults} from "vue";
import {uploadFile} from "@/plugin/cos/api";

interface Props {
  disabled?: boolean,
  accept?: string,
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  accept: "image/jpeg,image/jpg,image/png",
})

const emit = defineEmits(['change'])
const fileList = ref<UploadProps['fileList']>([]);
const customRequest = async (data: any) => {
  const url = await uploadFile(data)
  emit('change', url)
  if (Number(fileList.value?.length) > 0) {
    fileList.value?.pop()
  }
  fileList.value?.push({
    uid: data.file.uid,
    name: data.file.name,
    status: 'done',
    url: url,
  })
};
</script>

<style scoped>

</style>