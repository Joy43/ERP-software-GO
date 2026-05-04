Requisition → Approval → Purchase Order → Goods Receive → Inventory Update


📋 ধাপ ১: Requisition (চাহিদাপত্র)
কোনো কর্মী বা বিভাগ যখন কিছু কিনতে চায়, তখন সে একটা Requisition তৈরি করে। এটা মূলত একটা "আমার এটা দরকার" ফর্ম।
Requisition ৩ ধরনের:
ধরনমানেEmployee Requisitionএকজন কর্মী নিজের জন্য কিছু চাইছেDepartment Requisitionপুরো একটা বিভাগ কিছু চাইছেProject Requisitionকোনো প্রজেক্টের জন্য কিছু চাইছে
Requisition-এর Status গুলো:
DRAFT → PENDING → DEPARTMENT_APPROVED → FINANCE_APPROVED → APPROVED → ORDERED
                                                          ↘ REJECTED
                                                          ↘ CANCELLED

DRAFT = এখনো তৈরি হচ্ছে, submit হয়নি
PENDING = Submit করা হয়েছে, অনুমোদনের অপেক্ষায়
DEPARTMENT_APPROVED = বিভাগীয় প্রধান অনুমোদন দিয়েছেন
FINANCE_APPROVED = Finance বিভাগ অনুমোদন দিয়েছেন
APPROVED = সম্পূর্ণ অনুমোদিত
ORDERED = Purchase Order তৈরি হয়ে গেছে


✅ ধাপ ২: Approval (অনুমোদন প্রক্রিয়া)
Requisition submit করার পর সেটা ধাপে ধাপে অনুমোদন নিতে হয়। কে অনুমোদন দেবে সেটা Requisition-এর ধরনের উপর নির্ভর করে:
Employee Requisition  → Manager অনুমোদন → Finance/Admin অনুমোদন
Department Requisition → HOD (বিভাগীয় প্রধান) → Finance অনুমোদন  
Project Requisition   → Project Manager → Procurement অনুমোদন
এই অনুমোদনের ইতিহাস requisition_status_history টেবিলে সংরক্ষিত থাকে — কে কখন কী সিদ্ধান্ত নিল সব রেকর্ড থাকে।

🛒 ধাপ ৩: Purchase Order / PO (ক্রয় আদেশ)
Requisition অনুমোদিত হলে Purchase Order তৈরি হয়। এটা মূলত Supplier (সরবরাহকারী)-কে পাঠানো অর্ডার।
PO তৈরি হয় ২ ভাবে:

✅ Approved Requisition থেকে — স্বাভাবিক প্রক্রিয়া
❗ সরাসরি Manual — Requisition ছাড়াই সরাসরি PO করা যায়

PO-এর Status গুলো:
DRAFT → ISSUED → CONFIRMED → PARTIALLY_RECEIVED → FULLY_RECEIVED
                           ↘ CANCELLED

DRAFT = তৈরি হচ্ছে
ISSUED = Supplier-কে পাঠানো হয়েছে
CONFIRMED = Supplier নিশ্চিত করেছে
PARTIALLY_RECEIVED = কিছু পণ্য এসেছে, বাকি আসেনি
FULLY_RECEIVED = সব পণ্য পাওয়া গেছে

PO-তে কী কী তথ্য থাকে:

Supplier কে?
কোন পণ্য, কত quantity, কত দাম?
VAT, Discount হিসাব
Delivery date কখন?


📦 ধাপ ৪: Goods Receipt Note / GRN (পণ্য গ্রহণ)
Supplier পণ্য পাঠালে সেটা গ্রহণ করার সময় GRN তৈরি করা হয়। এটা মূলত "পণ্য পেয়েছি" এর রশিদ।
GRN ২ ধরনের:

✅ Against PO = PO-এর বিপরীতে পণ্য গ্রহণ (স্বাভাবিক)
❗ Direct Receive = PO ছাড়াই সরাসরি পণ্য গ্রহণ

GRN-এ কী তথ্য থাকে:

Challan নম্বর (Supplier-এর ডেলিভারি কাগজ)
Invoice নম্বর
কোন পণ্য কত quantity পাওয়া গেল
কত দামে কেনা হলো


📊 ধাপ ৫: Inventory Update (গুদামের হিসাব আপডেট)
GRN তৈরি হলে স্বয়ংক্রিয়ভাবে গুদামের stock বাড়ে।
location_stocks টেবিলে আপডেট হয়:
ক্ষেত্রমানেquantityবর্তমান মোট stocklast_costসর্বশেষ কত দামে কেনা হয়েছেaverage_costগড় ক্রয়মূল্যreserved_quantityRequisition-এ reserve আছে কত
stock_transactions টেবিলে প্রতিটি movement রেকর্ড থাকে:
GRN হলে     → quantity_change = +100 (বাড়ল)
Issue হলে   → quantity_change = -50  (কমল)
Transfer    → এক location থেকে অন্যটায় গেল
Adjustment  → ভুল ঠিক করা হলো

🗄️ Database টেবিলগুলোর সম্পর্ক সহজ ভাষায়
items (পণ্যের তালিকা)
  └── location_stocks (কোন গুদামে কত আছে)
  
requisitions (চাহিদাপত্র)
  └── requisition_items (কোন পণ্য কত চাই)
  └── requisition_status_history (অনুমোদনের ইতিহাস)

purchase_orders (ক্রয় আদেশ)
  └── purchase_order_items (কোন পণ্য কত অর্ডার)

goods_receipt_notes (পণ্য গ্রহণের রশিদ)
  └── grn_items (কোন পণ্য কত পাওয়া গেল)
  
stock_transactions (সকল stock এর movement এর লগ)

🔄 পুরো প্রক্রিয়া একটি উদাহরণ দিয়ে:

রহিম সাহেব অফিসের জন্য ১০টি চেয়ার কিনতে চান:


রহিম সাহেব Requisition তৈরি করলেন → Status: DRAFT
Submit করলেন → Status: PENDING
বিভাগীয় প্রধান অনুমোদন দিলেন → Status: DEPARTMENT_APPROVED
Finance অনুমোদন দিল → Status: APPROVED
Buyer একটি Purchase Order তৈরি করলেন Supplier-কে → PO Status: ISSUED
Supplier ১০টি চেয়ার দিল, GRN তৈরি হলো → PO Status: FULLY_RECEIVED
Inventory আপডেট হলো → গুদামে ১০টি চেয়ার যোগ হলো