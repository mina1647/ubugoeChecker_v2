<script lang="ts" setup>
import { ref } from 'vue'
// import UserMessage from '@/components/Message.vue'

const messages = ref<{ id: string; message: string }[]>([])
const query = ref<string>('')

const checker = async (name: string) => {
  const response = await fetch(`/api/ubugoe/${name}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  })
  if (!response.ok) {
    messages.value = [
      {
        id: 'error',
        message: '該当するtraQ IDが見つかりません',
      },
    ]
    return
  }

  const result = await response.json()
  messages.value = [] // 初期化

  for (let i = 0; i < result.length; i++) {
    const item = result[i]
    console.log('受信したメッセージ', item)

    setTimeout(() => {
      messages.value.push({
        id: item.userId, // ← traQのユーザーID（正しい）
        message: item.message, // ← メッセージ本文
      })
    }, i * 1000)
  }
}
</script>

<template>
  <div class="checker">
    <h1>うぶごえチェッカー</h1>
    <!-- <p>ここにうぶごえチェッカーの内容が入ります。</p> -->
    <v-text-field
      type="text"
      v-model="query"
      placeholder="traQIDを入力してください"
      class="nameinput"
    />
    <v-btn variant="outlined" @click="checker(query)">検索</v-btn>

    <!-- <div v-for="message in messages" :key="message.id">
      <UserMessage class="user" :id="message.id" :message="message.message" />
    </div> -->

    <div v-for="item in messages" :key="item.id" class="user">
      {{ item.message }}
    </div>
  </div>
</template>

<style scoped></style>
