<template>
  <div>
    <a-input v-model:value="username"></a-input>
    <a-input v-model:value="password"></a-input>
    <a-button @click="onRegister">点击注册</a-button>
    <a-button @click="onLogin">点击登录</a-button>
    <a-button @click="onMatch">点击匹配</a-button>
    <a-button @click="onWantLandlord">叫地主</a-button>
  </div>
</template>
<script lang="ts" setup>
import {ref} from 'vue';
import {message} from 'ant-design-vue';
import {Event_Type, Webhook_Type} from "@/plugin/api/model/lbddz";

let username = ref<string>("")
let password = ref<string>("")
let playerId = ref<string>("")
let gameId = ref<string>("")

//定义一个websocket
let websocket: WebSocket | null = null;

//初始化WebSocket
const wsInit = async () => {
  if (websocket) {
    console.log("websocket is available")
    return;
  }

  //判断当前浏览器是否支持WebSocket（固定写法）
  if ('WebSocket' in window) {
    websocket = new WebSocket('ws://127.0.0.1:8889/');
  } else {
    message.error('浏览器不支持websocket');
    return
  }

  //连接发生错误的回调方法
  websocket.onerror = function (e) {
    console.log('websocket 发生错误', e);
  };

  //连接成功建立的回调方法
  websocket.onopen = function (e) {
    console.log('websocket 建立连接', e);
  };

  //接收到消息的回调方法
  websocket.onmessage = function (e) {
    console.log('websocket 接收消息', e);
    if (e.data instanceof Blob) {
      const reader = new FileReader();
      reader.readAsText(e.data, "UTF-8");
      reader.onload = (_: any) => {
        if (typeof reader.result === "string") {
          const result = JSON.parse(reader.result);
          if (result.Webhook) {
            onHandler(result.Webhook)
          } else {
            onHandler(result)
          }
        }
      };
    } else {
      const result = JSON.parse(e.data);
      if (result.Webhook) {
        onHandler(result.Webhook)
      } else {
        onHandler(result)
      }
    }

  };

  //连接关闭的回调方法
  websocket.onclose = function (e) {
    console.log('websocket 关闭连接', e);
  };
};
wsInit()

// 关闭socket
const closeSocket = () => {
  websocket?.close();
};

// 发送消息
const sendMsg = (param: any) => {
  wsInit();
  console.log("send msg is ", JSON.stringify(param))
  websocket?.send(JSON.stringify(param));
}


// 点击发送事件
const clickSend = () => {
  //调用消息接口传消息
  try {
    const param = {};
    sendMsg(param)
  } catch (e) {
    console.log("err is ,", e)
  }
};

const onRegister = () => {
  sendMsg({
    Register: {
      username: username.value,
      password: password.value,
    }
  })
}

const onLogin = () => {
  sendMsg({
    Login: {
      username: username.value,
      password: password.value,
    }
  })
}

const onEvent = (data: lbddz.Event) => {
  sendMsg({
    Event: data
  })
}

const onMatch = () => {
  onEvent({
    type: Event_Type.TypeMatchPlayer,
    match_player: {
      player_id: playerId.value,
    },
  })
}

const onWantLandlord = () => {
  onEvent({
    type: Event_Type.TypeWantLandlord,
    want_landlord: {
      game_id: gameId.value,
      score: 1,
    }
  })
}

const onHandler = (data: lbddz.Webhook) => {
  console.log("handler 收到", data);
  switch (data.type) {
    case Number(Webhook_Type.TypeLoginResult):
      playerId.value = String(data.login?.player?.id)
      break
    case Number(Webhook_Type.TypeMatchResult):
      console.log("match is ,", data.match)
      break
    case Number(Webhook_Type.TypeGiveCard):
      data.give_card?.base_game?.gps?.forEach(function (value) {
        if (value.player_id === String(playerId.value)) {
          console.log("game is ,", data.give_card?.base_game?.g)
          console.log("giveCard is ,", value)
        }
      })
      break
    case Number(Webhook_Type.TypeStateChange):
      console.log("state change is ,", data.state_change)
      gameId.value = String(data.state_change?.base_game?.g?.id)
      break
    default:
      console.log("unknown type", data.type)
      break
  }
}
</script>

