{
  "average_cost": 145,
  "buyer_id": 17,
  "category_id": 1,
  "created_date": "2026-05-01",
  "current_stock": 5,
  "department_id": 1,
  "description": "High quality office chairs",
  "employee_id": 13,
  "expected_date": "2026-05-15",
  "inventory_type_id": 1,
  "item_id": 5,
  "item_type_id": 2,
  "last_cost": 150.5,
  "location_id": 9,
  "minor_category_id": 2,
  "office_id": 1,
  "project_id": 1,
  "remarks": "Office supplies needed for Q2",
  "request_quantity": 20,
  "requisition_type": "EMPLOYEE",
  "sub_category_id": 1,
  "supplier_id": 4
}


//-------- requasition -------
{
  "buyer_id": 17,
  "department_id": 1,
  "description": "string",
  "employee_id": 13,
  "expected_date": "2026-05-15",
  "inventory_type_id": 1,
  "items": [
    {
      "average_cost": 30,
      "category_id": 1,
      "current_stock": 20,
      "description": "string",
      "item_id": 5,
      "item_type_id": 2,
      "last_cost": 150,
      "minor_category_id": 2,
      "request_quantity": 20,
      "sub_category_id": 1
    }
  ],
  "location_id": 9,
  "office_id": 1,
  "project_id": 1,
  "remarks": "string welcome",
  "requisition_type": "EMPLOYEE",
  "supplier_id": 4
}

// --parchase order--


{
"po_date": "2026-05-02",
"order_type": "REQUISITION_BASED",
"requisition_id": 4,
"office_id": 1,
"location_id": 9,
"supplier_id": 4
}

//----- grn payload----

{
  "challan_date": "2026-05-01",
  "challan_no": "wonder to make us",
  "delivery_number": "2026-05-15",
  "file_id": 1,
  "grn_date": "2026-05-15",
  "items": [
    {
      "category_id":1,
      "discount_percentage": 10,
      "item_id": 5,
      "minor_category_id": 2,
      "po_item_id": 2,
      "purchase_price": 60,
      "received_quantity": 3,
      "remarks": "best price string",
      "sub_category_id": 1,
      "uom_id": 1,
      "vat_percentage": 40
    }
  ],
  "location_id": 9,
  "office_id": 1,
  "payment_method_id": 1,
  "po_id": 3,
  "receive_type": "DIRECT",
  "received_by_id": 17,
  "remarks": "string",
  "requisition_id": 7,
  "sales_invoice_number": "string",
  "shipment_document_number": "string",
  "shipping_address": "string",
  "supplier_id": 4,
  "vat_challan_number": "4747string"
}

// ----------- advance payments -----------------

{
  "account_head_id": 17,
  "amount": 4080,
  "cash_amount": 600,
  "lc_no": "5776string",
  "narration": "5757string",
  "office_id": 1,
  "payment_date": "2026-08-01",
  "payment_mode_id": 1,
  "po_id": 2,
  "supplier_head_id": 17
}

------------------  payment by grn -------------------
{
  "adjustment_amount": 700,
  "grn_id": 2,
  "money_receipt_no": "7756string",
  "office_head_id": 17,
  "office_id": 1,
  "payable_amount": 650,
  "paying_amount": 40,
  "payment_date": "2026-08-01",
  "payment_mode_id": 1,
  "supplier_id": 4
}

--------- supplier bills-----------
