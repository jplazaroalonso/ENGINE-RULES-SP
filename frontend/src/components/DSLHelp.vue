<template>
  <q-dialog v-model="showDialog" persistent>
    <q-card class="dsl-help-dialog" style="min-width: 800px; max-width: 1000px;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">DSL (Domain Specific Language) Help</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-tabs v-model="activeTab" class="text-primary">
          <q-tab name="overview" label="Overview" />
          <q-tab name="syntax" label="Syntax" />
          <q-tab name="examples" label="Examples" />
          <q-tab name="functions" label="Functions" />
        </q-tabs>

        <q-tab-panels v-model="activeTab" animated>
          <!-- Overview Tab -->
          <q-tab-panel name="overview">
            <div class="dsl-help-content">
              <h3>What is DSL?</h3>
              <p>
                The Domain Specific Language (DSL) is a simple, readable language designed specifically 
                for defining business rules. It allows you to express complex business logic in a 
                natural, human-readable format.
              </p>

              <h3>Basic Structure</h3>
              <p>Every DSL rule follows this basic structure:</p>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF condition THEN action</code></pre>
                </q-card-section>
              </q-card>

              <h3>Key Components</h3>
              <ul>
                <li><strong>IF</strong> - Defines the condition that must be met</li>
                <li><strong>THEN</strong> - Defines the action to take when the condition is true</li>
                <li><strong>Variables</strong> - Reference to data (e.g., customer.age, order.total)</li>
                <li><strong>Operators</strong> - Comparison and logical operators</li>
                <li><strong>Functions</strong> - Predefined actions (e.g., apply_discount, send_notification)</li>
              </ul>
            </div>
          </q-tab-panel>

          <!-- Syntax Tab -->
          <q-tab-panel name="syntax">
            <div class="dsl-help-content">
              <h3>Condition Syntax</h3>
              <p>Conditions can use various operators and data types:</p>

              <h4>Comparison Operators</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>customer.age > 18
order.total >= 100
product.category == "electronics"
customer.status != "inactive"</code></pre>
                </q-card-section>
              </q-card>

              <h4>Logical Operators</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>customer.age > 18 AND customer.status == "active"
order.total > 50 OR customer.is_vip == true
NOT customer.is_blocked</code></pre>
                </q-card-section>
              </q-card>

              <h4>Data Types</h4>
              <ul>
                <li><strong>Numbers:</strong> 18, 100.50, -25</li>
                <li><strong>Strings:</strong> "electronics", "active", "premium"</li>
                <li><strong>Booleans:</strong> true, false</li>
                <li><strong>Variables:</strong> customer.age, order.total, product.category</li>
              </ul>

              <h3>Action Syntax</h3>
              <p>Actions are predefined functions that perform specific operations:</p>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>apply_discount(10)
send_notification("welcome")
set_priority("high")
add_tag("vip")</code></pre>
                </q-card-section>
              </q-card>
            </div>
          </q-tab-panel>

          <!-- Examples Tab -->
          <q-tab-panel name="examples">
            <div class="dsl-help-content">
              <h3>Common Rule Examples</h3>

              <h4>Customer Age Discount</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF customer.age > 65 THEN apply_discount(10)</code></pre>
                </q-card-section>
              </q-card>
              <p>Applies a 10% discount to customers over 65 years old.</p>

              <h4>High Value Order</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF order.total > 100 THEN set_priority("high")</code></pre>
                </q-card-section>
              </q-card>
              <p>Sets high priority for orders over $100.</p>

              <h4>VIP Customer Benefits</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF customer.is_vip == true AND order.total > 50 THEN apply_discount(15)</code></pre>
                </q-card-section>
              </q-card>
              <p>Applies 15% discount to VIP customers with orders over $50.</p>

              <h4>Product Category Rules</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF product.category == "electronics" AND order.total > 200 THEN apply_discount(5)</code></pre>
                </q-card-section>
              </q-card>
              <p>Applies 5% discount to electronics orders over $200.</p>

              <h4>Time-based Rules</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF order.time.hour >= 9 AND order.time.hour <= 17 THEN set_priority("normal")</code></pre>
                </q-card-section>
              </q-card>
              <p>Sets normal priority for orders placed during business hours.</p>

              <h4>Complex Conditions</h4>
              <q-card flat bordered class="code-example">
                <q-card-section>
                  <pre><code>IF (customer.age > 18 AND customer.status == "active") OR customer.is_vip == true THEN apply_discount(20)</code></pre>
                </q-card-section>
              </q-card>
              <p>Applies 20% discount to active customers over 18 or VIP customers.</p>
            </div>
          </q-tab-panel>

          <!-- Functions Tab -->
          <q-tab-panel name="functions">
            <div class="dsl-help-content">
              <h3>Available Functions</h3>

              <h4>Discount Functions</h4>
              <q-card flat bordered class="function-card">
                <q-card-section>
                  <h5>apply_discount(percentage)</h5>
                  <p>Applies a percentage discount to the order.</p>
                  <p><strong>Parameters:</strong> percentage (number) - The discount percentage (0-100)</p>
                  <p><strong>Example:</strong> <code>apply_discount(15)</code> - Applies 15% discount</p>
                </q-card-section>
              </q-card>

              <h4>Priority Functions</h4>
              <q-card flat bordered class="function-card">
                <q-card-section>
                  <h5>set_priority(level)</h5>
                  <p>Sets the priority level for the order or customer.</p>
                  <p><strong>Parameters:</strong> level (string) - "low", "normal", "high", "urgent"</p>
                  <p><strong>Example:</strong> <code>set_priority("high")</code> - Sets high priority</p>
                </q-card-section>
              </q-card>

              <h4>Notification Functions</h4>
              <q-card flat bordered class="function-card">
                <q-card-section>
                  <h5>send_notification(message)</h5>
                  <p>Sends a notification to the customer.</p>
                  <p><strong>Parameters:</strong> message (string) - The notification message</p>
                  <p><strong>Example:</strong> <code>send_notification("Welcome to our store!")</code></p>
                </q-card-section>
              </q-card>

              <h4>Tagging Functions</h4>
              <q-card flat bordered class="function-card">
                <q-card-section>
                  <h5>add_tag(tag)</h5>
                  <p>Adds a tag to the customer or order.</p>
                  <p><strong>Parameters:</strong> tag (string) - The tag to add</p>
                  <p><strong>Example:</strong> <code>add_tag("vip")</code> - Adds VIP tag</p>
                </q-card-section>
              </q-card>

              <h4>Status Functions</h4>
              <q-card flat bordered class="function-card">
                <q-card-section>
                  <h5>set_status(status)</h5>
                  <p>Sets the status of the customer or order.</p>
                  <p><strong>Parameters:</strong> status (string) - The status value</p>
                  <p><strong>Example:</strong> <code>set_status("approved")</code> - Sets approved status</p>
                </q-card-section>
              </q-card>

              <h3>Available Variables</h3>
              <p>You can reference the following data in your conditions:</p>
              <ul>
                <li><strong>customer.*</strong> - Customer information (age, status, is_vip, etc.)</li>
                <li><strong>order.*</strong> - Order information (total, time, items, etc.)</li>
                <li><strong>product.*</strong> - Product information (category, price, stock, etc.)</li>
                <li><strong>location.*</strong> - Location information (country, city, region, etc.)</li>
              </ul>
            </div>
          </q-tab-panel>
        </q-tab-panels>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Close" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const activeTab = ref('overview')

const showDialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

<style lang="scss" scoped>
.dsl-help-dialog {
  max-height: 80vh;
}

.dsl-help-content {
  max-height: 60vh;
  overflow-y: auto;
  
  h3 {
    color: var(--carrefour-blue);
    margin-top: 24px;
    margin-bottom: 12px;
  }
  
  h4 {
    color: var(--carrefour-dark-blue);
    margin-top: 20px;
    margin-bottom: 8px;
  }
  
  h5 {
    color: var(--carrefour-blue);
    margin: 0 0 8px 0;
    font-family: 'Courier New', monospace;
  }
  
  p {
    margin-bottom: 12px;
    line-height: 1.6;
  }
  
  ul {
    margin-bottom: 16px;
    
    li {
      margin-bottom: 4px;
    }
  }
}

.code-example {
  background-color: var(--carrefour-gray-50);
  border-left: 4px solid var(--carrefour-blue);
  margin: 12px 0;
  
  pre {
    margin: 0;
    font-family: 'Courier New', monospace;
    font-size: 14px;
    color: var(--carrefour-dark-blue);
  }
}

.function-card {
  margin: 12px 0;
  border-left: 4px solid var(--carrefour-green);
  
  p {
    margin-bottom: 8px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  code {
    background-color: var(--carrefour-gray-100);
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    color: var(--carrefour-dark-blue);
  }
}
</style>
