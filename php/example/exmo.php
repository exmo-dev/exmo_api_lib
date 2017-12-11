<?php

    require ('api.php');
    
    class exmo extends api {
        public function trades ($pair) {
            return $this -> query ('trades', [
                'pair' => implode (',', $pair)
            ]);
        }
        public function order_book ($pair, $limit = 100) {
            return $this -> query ('order_book', [
                'pair' => implode (',', $pair),
                'limit' => $limit
            ]);
        }
        public function ticker () {
            return $this -> query ('ticker');
        }
        public function pair_settings () {
            return $this -> query ('pair_settings');
        }
        public function currency () {
            return $this -> query ('currency');
        }
        public function user_info () {
            return $this -> query ('user_info');
        }
        public function order_create ($pair, $quantity, $price, $type) {
            return $this -> query ('order_create', [
                'pair' => $pair,
                'quantity' => $quantity,
                'price' => $price,
                'type' => $type
            ]);
        }
        public function order_cancel ($order_id) {
            return $this -> query ('order_cancel', [
                'order_id' => $order_id
            ]);
        }
        public function user_open_orders () {
            return $this -> query ('user_open_orders');
        }
        public function user_trades ($pair, $offset = 0, $limit = 100) {
            return $this -> query ('user_trades', [
                'pair' => $pair,
                'offset' => $offset,
                'limit' => $limit
            ]);
        }
        public function user_cancelled_orders ($offset, $limit = 100) {
            return $this -> query ('user_cancelled_orders', [
                'offset' => $offset,
                'limit' => $limit
            ]);
        }
        public function order_trades ($order_id) {
            return $this -> query ('order_trades', [
                'order_id' => $order_id
            ]);
        }
        public function required_amount ($pair, $quantity) {
            return $this -> query ('required_amount', [
                'pair' => $pair,
                'quantity' => $quantity
            ]);
        }
        public function deposit_address () {
            return $this -> query ('deposit_address');
        }
        public function withdraw_crypt ($amount, $currency, $address, $invoice) {
            return $this -> query ('withdraw_crypt', [
                'amount' => $amount,
                'currency' => $currency,
                'address' => $address,
                'invoice' => $invoice
            ]);
        }
        public function withdraw_get_txid ($task_id) {
            return $this -> query ('withdraw_get_txid', [
                'task_id' => $task_id
            ]);
        }
        public function excode_create ($currency, $amount) {
            return $this -> query ('excode_create', [
                'currency' => $currency,
                'amount' => $amount
            ]);
        }
        public function excode_load ($code) {
            return $this -> query ('excode_load', [
                'code' => $code
            ]);
        }
        public function wallet_history ($date = '') {
            return $this -> query ('wallet_history', [
                'date' => $date
            ]);
        }
    }
