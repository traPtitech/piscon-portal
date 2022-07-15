<template>
  <div class="logo">
    PISCON
  </div>
</template>
<script>
import { useColors } from "vuestic-ui";

export default {
  name: "VaIconVuestic",

  props: {
    height: { type: [Number, String], default: 16 },
    color: { type: [String], default: "primary" },
  },

  computed: {
    colorsComputed() {
      const { getColor, shiftHSLAColor } = useColors();
      const color = getColor(this.color, "primary");

      return { start: color, end: shiftHSLAColor(color, { l: -20 }) };
    },
  },
};
</script>

<style lang="scss" scoped>
.logo {
  font-size: 1.5rem;
  color: var(--va-primary);
}
</style>