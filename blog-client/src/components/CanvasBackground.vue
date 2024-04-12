<template>
  <div class="star_content">
    <canvas id="star_canvas" ref="star_canvas"></canvas>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import { Star, Move, Draw, DrawLine, generateStar } from "./index";

const star_canvas = ref<HTMLCanvasElement>();
let ctx: CanvasRenderingContext2D;

let timer: any;

const initStar = () => {
  ctx = star_canvas.value!.getContext("2d") as CanvasRenderingContext2D;
  // ctx.fillStyle = "white";
  // ctx.strokeStyle = "white";
  ctx.fillStyle = "rgba(45,140,210,0.2)";
  ctx.strokeStyle = "rgba(45,140,210,0.2)";
  let starArr: Star[] = [];
  for (let i = 0; i < 100; i++) {
    let star = generateStar();
    starArr.push(star);
  }

  // 鼠标star
  const mouse_star = generateStar();
  mouse_star.speedX = 0;
  mouse_star.speedY = 0;
  // 鼠标star跟随移动
  star_canvas.value!.onmousemove = e => {
    mouse_star.x = e.offsetX;
    mouse_star.y = e.offsetY;
  };

  timer = setInterval(() => {
    ctx.clearRect(0, 0, document.documentElement.clientWidth, document.documentElement.clientHeight);
    // 鼠标star移动
    Draw(ctx, mouse_star);
    // star移动
    starArr.forEach((star: Star) => {
      Draw(ctx, star);
      Move(star);
    });

    // 比较star和所有其他star的距离，小于50连线
    starArr.forEach((star: Star, index: number) => {
      for (let i = index + 1; i < starArr.length; i++) {
        if (Math.abs(star.x - starArr[i].x) < 150 && Math.abs(star.y - starArr[i].y) < 150) {
          DrawLine(star.x, star.y, starArr[i].x, starArr[i].y, ctx);
        }
      }
    });
    // 比较鼠标star和所有star的距离
    for (let i = 0; i < starArr.length; i++) {
      if (Math.abs(mouse_star.x - starArr[i].x) < 100 && Math.abs(mouse_star.y - starArr[i].y) < 100) {
        DrawLine(mouse_star.x, mouse_star.y, starArr[i].x, starArr[i].y, ctx);
      }
    }
  }, 50);

  // 点击添加star
  star_canvas.value!.onclick = e => {
    for (let i = 0; i < 5; i++) {
      let star = generateStar();
      star.x = e.offsetX;
      star.y = e.offsetY;
      starArr.push(star);
    }
  };
};
const initCanvas = () => {
  ctx = star_canvas.value?.getContext("2d") as CanvasRenderingContext2D;
  star_canvas.value!.width = document.documentElement.clientWidth - 18
  star_canvas.value!.height = document.documentElement.clientHeight - 8
  initStar();
};
onMounted(() => {
  initCanvas();
});
onBeforeUnmount(() => {
  clearInterval(timer);
});
</script>

<style lang="css" scoped>
#star_canvas {
  background-color: #333;
}
</style>