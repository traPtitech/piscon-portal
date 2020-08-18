<template>
    <transition name="custom-modal">
    <div class="custom-modal-mask">
        <div class="custom-modal-wrapper">
        <div class="custom-modal-container">

            <div class="custom-modal-header">
              <slot name="header">
                  
              </slot>
            </div>

            <div class="my-custom-modal-header">
              <slot name="body">
                  default body
              </slot>
            </div>

            <div class="custom-modal-footer">
            <slot name="footer">
              <div class="custom-modal-default-button">
                <button class="btn btn-micro btn-info"  @click="$emit('operation')">
                はい
                </button>
                <button class="btn btn-micro btn-light" @click="$emit('close')">
                やめる
                </button>
              </div >
            </slot>
            </div>
        </div>
        </div>
    </div>
    </transition>
</template>

<script>
      // register custom-modal component
export default {
  name: 'customModal'
}
</script>

<style>
.custom-modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: table;
  transition: opacity 0.3s ease;
}

.custom-modal-wrapper {
  display: table-cell;
  vertical-align: middle;
  /* position: relative; */
}

.custom-modal-container {
  width: 300px;
  height:300px;
  margin: 0px auto;
  padding: 20px 30px;
  background-color: #fff;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
  font-family: Helvetica, Arial, sans-serif;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  align-content: center;
}

.custom-modal-header h3 {
  margin-top: 0;
  color: #42b983;
  /* flex-direction: row; */
}

.my-custom-modal-body  {
  /* margin: 20px 0; */
  padding: 0;
  /* justify-content: center; */
}
.custom-modal-footer {
  /* position: absolute; */
  /* top: 300px; */
  /* width:100%; */
  padding: 0;
  /* justify-content: flex-end; */

}

.custom-modal-default-button {
  /* float: right; */
  /* bottom: 10px */
  /* position: absolute; */
  /* top: 300px; */

  margin: 0px auto;

}

/*
 * The following styles are auto-applied to elements with
 * transition="custom-modal" when their visibility is toggled
 * by Vue.js.
 *
 * You can easily play with the custom-modal transition by editing
 * these styles.
 */

.custom-modal-enter {
  opacity: 0;
}

.custom-modal-leave-active {
  opacity: 0;
}

/* .custom-modal-enter .custom-modal-container,
.custom-modal-leave-active .custom-modal-container {
  -webkit-transform: scale(1.1);
  transform: scale(1.1);
} */
</style>