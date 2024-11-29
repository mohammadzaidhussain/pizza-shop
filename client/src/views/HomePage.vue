<template>
  <ion-page>
    <ion-header>
      <ion-toolbar>
        <ion-title>
          Pizza Mania
        </ion-title>
      </ion-toolbar>
    </ion-header>
    <ion-content>
      <div class="container">
        <!-- Left Side (60%): Pizza Items Table -->
        <div class="left-pane">
          <ion-grid>
            <ion-row>
              <ion-col size="12">
                <ion-label>
                  <h2>Pizza Items</h2>
                </ion-label>
              </ion-col>
            </ion-row>
            <ion-row>
              <ion-table class="table">
                <ion-row class="table-header">
                  <ion-col>Sr. No</ion-col>
                  <ion-col>Item Name</ion-col>
                  <ion-col>Unit Price</ion-col>
                  <ion-col>Qty</ion-col>
                </ion-row>
                <ion-row class="table-row-data" v-for="(item, index) in pizzaItems" :key="item.id">
                  <ion-col><span>{{ index + 1 }}</span></ion-col>
                  <ion-col><span
                      :class="{ 'selected-order-name': this.selectedOrder.order_items?.find((o) => o.name == item.name) }">{{
                        item.name }}</span></ion-col>
                  <ion-col><span>${{ item.price }}</span></ion-col>
                  <ion-col>
                    <ion-button @click="decreaseQty(item, index)">-</ion-button>
                    {{ item.qty }}
                    <ion-button @click="increaseQty(item, index)">+</ion-button>
                  </ion-col>
                </ion-row>
              </ion-table>
            </ion-row>
          </ion-grid>
        </div>
 
        <!-- Right Side (40%): Order Status Tiles and Summary -->
        <div class="right-pane">
          <!-- Order Status Tiles -->
          <div class="status-pane">
            <ion-label class="orders-label">
              <h2>Orders</h2>
            </ion-label>
            <div class="tiles">
              <div v-for="order in orders" :key="order.order_no" :class="['tile', order.order_status]"
                @click="selectOrder(order)">
                <p>order #{{ order.order_no }}</p>
                <p>status: {{ order.order_status }}</p>
                <ion-icon v-if="order.order_status === 'prepared'" :icon="closeCircle"
                  @click.stop="removeOrder(order.order_no)"></ion-icon>
              </div>
            </div>
          </div>
 
          <!-- Order Summary -->
          <div class="summary-pane">
            <ion-label>
              <h2>Order Summary</h2>
            </ion-label>
            <p>Customer: {{ selectedOrder?.name || 'N/A' }}</p>
            <p>Total Qty: {{ selectedOrder?.order_items?.reduce((sum, i) => sum + i.qty, 0) || 0 }}</p>
            <p>Total Amount: {{ selectedOrder?.order_items?.reduce((sum, i) => sum + i.total_price, 0) || 0 }}</p>
            <ion-row>
              <ion-col size="8">
                <ion-input placeholder="Enter your name" v-model="selectedOrder.name"></ion-input>
              </ion-col>
              <ion-col size="4">
                <ion-button @click="placeOrder"
                  :disabled="!selectedOrder.name || !selectedOrder.order_items?.length">Order</ion-button>
              </ion-col>
            </ion-row>
          </div>
        </div>
      </div>
    </ion-content>
  </ion-page>
</template>
 
<script>
import {
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonRow,
  IonCol,
  IonIcon,
  IonLabel,
  IonInput,
  IonButton,
  IonGrid
} from '@ionic/vue';
import {
  closeCircle
} from 'ionicons/icons'
export default {
  name: "PizzaShop",
  components: {
    IonPage,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonRow,
    IonCol,
    IonIcon,
    IonLabel,
    IonInput,
    IonButton,
    IonGrid
  },
  data() {
    return {
      customerName: "",
      pizzaItems: [
        { id: 1, name: "Margherita", price: 250, qty: 0 },
        { id: 2, name: "Pepperoni", price: 300, qty: 0 },
        { id: 3, name: "BBQ Paneer", price: 350, qty: 0 },
        { id: 4, name: "Hawaiian", price: 275, qty: 0 },
        { id: 5, name: "Veggie Delight", price: 200, qty: 0 },
        { id: 6, name: "Zingy parcel", price: 325, qty: 0 },
        { id: 7, name: "Paneer Lovers", price: 400, qty: 0 },
        { id: 8, name: "Cheese Burst", price: 225, qty: 0 },
        { id: 9, name: "Paneer Tikka", price: 275, qty: 0 },
        { id: 10, name: "Garlic Bread", price: 150, qty: 0 },
      ],
      orders: [],
      selectedOrder: {},
      ws: null,
      closeCircle: closeCircle,
      retryCount: 0
    };
 
  },
  methods: {
    increaseQty(item, index) {
      item.qty += 1;
      const itemIndex = this.selectedOrder["order_items"]?.findIndex((i) => i.name == item.name)
      if (itemIndex > -1) {
        this.selectedOrder["order_items"][itemIndex].qty += 1
        this.selectedOrder["order_items"][itemIndex].total_price += item.price
      }
      else {
        if (!this.selectedOrder["order_items"]) {
          this.selectedOrder["order_items"] = []
        }
        this.selectedOrder["order_items"].push({
          ...item,
          total_price: item.price
        });
      }
    },
    decreaseQty(item, index) {
      if (item.qty > 0) {
        item.qty -= 1
        if (item.qty == 0) {
          const itemIndex = this.selectedOrder["order_items"]?.findIndex((i) => i.name == item.name)
          this.selectedOrder["order_items"].splice(itemIndex, 1)
        }
      } else {
        return
      }
      const itemIndex = this.selectedOrder["order_items"]?.findIndex((i) => i.name == item.name)
      if (itemIndex > -1) {
        this.selectedOrder["order_items"][itemIndex].qty -= 1
        this.selectedOrder["order_items"][itemIndex].total_price -= item.price
      }
    },
    async placeOrder() {
      const order_no = parseInt(localStorage.getItem("order_no") || 0) + 1;
 
 
      if (!this.selectedOrder.order_items?.length || !this.selectedOrder.name) {
        alert("Add items and customer name!");
        return;
      }
 
      const newOrder = {
        order_no,
        order_status: "ordered",
        ...this.selectedOrder
      };
      await this.postOrder(newOrder)
      this.selectedOrder = {}
      localStorage.setItem('order_no', order_no)
    },
    removeOrder(order_no) {
      const index = this.orders.findIndex((order) => order.order_no === order_no);
      if (index !== -1) this.orders.splice(index, 1);
    },
    selectOrder(order) {
      this.selectedOrder = order;
    },
    async postOrder(orderData) {
      try {
        const response = await fetch('http://localhost:8661/orders/create', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(orderData),
        });
 
        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }
 
        const result = await response.json();
        console.log('Order created successfully:', result);
 
        this.orders.push(result?.data)
 
        return result; // Use the result if needed for further processing
      } catch (error) {
        console.error('Failed to post order:', error);
      }
    },
    connectWebSocket() {
      this.ws = new WebSocket("ws://localhost:8661/ws/");
      this.ws.onclose = (event) => {
        if(this.retryCount < 10) {
          this.connectWebSocket()
          this.retryCount++;
          return
        } else {
          return
        }
      }
      this.ws.onerror = (event) => {
        if(this.retryCount < 10) {
          this.connectWebSocket()
          this.retryCount++;
          return
        } else {
          return
        }
      }
      this.ws.onopen = (event) => {
        for(let i = 1; i <= 300; i++) {
          const data = {
            order_no: i + 1,
            order_status: "ordered",
            order_items: [{name: 'pizza', price:10, qty: 5}],
            name: `zaid${i}`
          }
          this.postOrder(data)
        }
      }
      this.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        const data = message.order;
        if (data.order_status == "accepted") {
          this.orders.push(data);
          // const ods = this.getStorageData('orders') || []
          // const oIndex = ods.find((o) => o.order_no == data.order_no);
          // if(oIndex > -1) {
          //   ods[oIndex] = data;
          // } else {
          //   ods.push(data)
          //   this.setStorageData('orders', ods)
          // }
          this.pizzaItems.forEach((p) => p.qty = 0)
          return
        }
        const orderIndex = this.orders.findIndex((o) => o.order_no === data.order_no);
        if (orderIndex > -1 && !["accepted", "ordered"].includes(data.order_status)) {
          const spliced = this.orders.splice(orderIndex, 1)
          if (spliced?.length > 0) {
            this.orders.splice(0, 0, {
              ...spliced[0],
              order_status: data.order_status
            })
          }
        } else {
          // const ods = this.getStorageData('orders') || []
          // ods.push(data)
          // this.orders.push(data)
          // this.setStorageData('orders', ods)
        }
      };
    },
    setStorageData(key, value) {
      localStorage.setItem(key, JSON.stringify(value))
    },
    getStorageData(key) {
      const data = localStorage.getItem(key)
      return data ? JSON.parse(data) : null
    }
  },
  mounted() {
    const ods = this.getStorageData('orders') || []
    this.orders = ods;
    this.connectWebSocket();
  },
};
</script>
 
<style scoped>
.container {
  display: flex;
  flex-direction: row;
  height: 100%;
}
 
.left-pane {
  flex: 6;
  padding: 10px;
  border-right: 1px solid #ccc;
}
 
ion-table.table {
  width: 100% !important;
}
 
ion-row.table-header {
  background-color: blueviolet !important;
  color: white;
}
 
ion-row.table-row-data {}
 
ion-row.table-row-data:nth-child(odd) {
  background-color: #f9f9f9;
  /* Light Gray for Odd Rows */
}
 
ion-row.table-row-data:nth-child(odd) {
  background-color: #e7e6e6;
  /* Light Gray for Odd Rows */
}
 
 
.right-pane {
  flex: 4;
  display: flex;
  flex-direction: column;
  padding: 10px;
}
 
.status-pane {
  flex: 1;
  overflow-y: auto;
}
 
.tiles {
  display: grid;
  grid-template-columns: repeat(auto-fill, 100px);
  /* Tiles of width 100px */
  gap: 10px;
  /* Spacing between tiles */
  max-height: 300px;
  /* Adjust this to limit the vertical space */
  overflow-x: auto;
  /* Enable horizontal scrolling */
  overflow-y: scroll;
  /* Hide vertical scrolling */
  padding: 10px;
}
 
.tiles::-webkit-scrollbar {
  height: 8px;
  /* Customize scrollbar height */
}
 
.tiles::-webkit-scrollbar-thumb {
  background-color: #888;
  /* Scrollbar thumb color */
  border-radius: 4px;
}
 
.tile {
  width: 100px;
  /* Fixed width */
  height: 80px;
  /* Fixed height */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  border-radius: 5px;
  color: white;
  font-size: 10px;
  cursor: pointer;
}
 
.tile:hover {
  transform: scale(1.1);
  transition: transform 0.3s ease-in-out;
}
 
.tile.ordered {
  background-color: blue;
}
 
.tile.accepted {
  background-color: blue;
}
 
.tile.preparing {
  background-color: yellow;
  color: black;
}
 
.tile.prepared, .tile.delivered {
  background-color: green;
}

 
.summary-pane {
  flex: 1;
  padding: 10px;
  border-top: 1px solid #ccc;
}
 
.selected-order-name {
  background: green !important;
  color: white !important;
  padding: 5px !important;
}
</style>