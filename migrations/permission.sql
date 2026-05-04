INSERT INTO `assmi_super_shop`.`permissions`
(`group_name`, `name`, `slug`, `description`, `created_at`, `updated_at`)
VALUES

-- ================= USERS =================
('users', 'View Users', 'users.view', 'Can view user list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('users', 'Create User', 'users.create', 'Can create new user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('users', 'Edit User', 'users.edit', 'Can edit user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('users', 'Delete User', 'users.delete', 'Can delete user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= ROLES =================
('roles', 'View Roles', 'roles.view', 'Can view roles list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('roles', 'Create Role', 'roles.create', 'Can create new role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('roles', 'Edit Role', 'roles.edit', 'Can edit role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('roles', 'Delete Role', 'roles.delete', 'Can delete role', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= DESIGNATIONS =================
('designations', 'View Designations', 'designations.view', 'Can view designation list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('designations', 'Create Designation', 'designations.create', 'Can create designation', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('designations', 'Edit Designation', 'designations.edit', 'Can edit designation', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('designations', 'Delete Designation', 'designations.delete', 'Can delete designation', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= DEPARTMENTS =================
('departments', 'View Departments', 'departments.view', 'Can view department list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('departments', 'Create Department', 'departments.create', 'Can create department', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('departments', 'Edit Department', 'departments.edit', 'Can edit department', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('departments', 'Delete Department', 'departments.delete', 'Can delete department', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= OFFICES =================
('offices', 'View Offices', 'offices.view', 'Can view office list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('offices', 'Create Office', 'offices.create', 'Can create office', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('offices', 'Edit Office', 'offices.edit', 'Can edit office', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('offices', 'Delete Office', 'offices.delete', 'Can delete office', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= PAYMENT MODES =================
('payment_modes', 'View Payment Modes', 'payment_modes.view', 'Can view payment modes', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= RESPONSIBILITY TRANSFERS =================
('responsibility_transfers', 'View Responsibility Transfers', 'responsibility_transfers.view', 'Can view responsibility transfers', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('responsibility_transfers', 'Create Responsibility Transfer', 'responsibility_transfers.create', 'Can create responsibility transfer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('responsibility_transfers', 'Approve Responsibility Transfer', 'responsibility_transfers.approve', 'Can approve responsibility transfer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('responsibility_transfers', 'Delete Responsibility Transfer', 'responsibility_transfers.delete', 'Can delete responsibility transfer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= TRANSACTION DELETE HISTORY =================
('transaction_history', 'View Transaction History', 'transaction_history.view', 'Can view transaction history', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('transaction_history', 'Delete Transaction History', 'transaction_history.delete', 'Can delete transaction history', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= SUPPLIERS =================
('suppliers', 'View Suppliers', 'suppliers.view', 'Can view supplier list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('suppliers', 'Create Supplier', 'suppliers.create', 'Can create new supplier', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('suppliers', 'Edit Supplier', 'suppliers.edit', 'Can edit supplier', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('suppliers', 'Delete Supplier', 'suppliers.delete', 'Can delete supplier', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= PARTNER GROUPS =================
('partner_groups', 'View Partner Groups', 'partner_groups.view', 'Can view partner group list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_groups', 'Create Partner Group', 'partner_groups.create', 'Can create partner group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_groups', 'Edit Partner Group', 'partner_groups.edit', 'Can edit partner group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_groups', 'Delete Partner Group', 'partner_groups.delete', 'Can delete partner group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= PARTNER SUB GROUPS =================
('partner_sub_groups', 'View Partner Sub Groups', 'partner_sub_groups.view', 'Can view partner sub group list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_sub_groups', 'Create Partner Sub Group', 'partner_sub_groups.create', 'Can create partner sub group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_sub_groups', 'Edit Partner Sub Group', 'partner_sub_groups.edit', 'Can edit partner sub group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('partner_sub_groups', 'Delete Partner Sub Group', 'partner_sub_groups.delete', 'Can delete partner sub group', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= DISTRICTS =================
('districts', 'View Districts', 'districts.view', 'Can view district list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('districts', 'Create District', 'districts.create', 'Can create district', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('districts', 'Edit District', 'districts.edit', 'Can edit district', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('districts', 'Delete District', 'districts.delete', 'Can delete district', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= THANAS =================
('thanas', 'View Thanas', 'thanas.view', 'Can view thana list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('thanas', 'Create Thana', 'thanas.create', 'Can create thana', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('thanas', 'Edit Thana', 'thanas.edit', 'Can edit thana', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('thanas', 'Delete Thana', 'thanas.delete', 'Can delete thana', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= TAX BRACKETS =================
('tax_brackets', 'View Tax Brackets', 'tax_brackets.view', 'Can view tax bracket list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tax_brackets', 'Create Tax Bracket', 'tax_brackets.create', 'Can create tax bracket', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tax_brackets', 'Edit Tax Bracket', 'tax_brackets.edit', 'Can edit tax bracket', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tax_brackets', 'Delete Tax Bracket', 'tax_brackets.delete', 'Can delete tax bracket', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


-- ================= CUSTOMERS =================
('customers', 'View Customers', 'customers.view', 'Can view customer list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('customers', 'Create Customer', 'customers.create', 'Can create new customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('customers', 'Edit Customer', 'customers.edit', 'Can edit customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('customers', 'Delete Customer', 'customers.delete', 'Can delete customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--     iteam profile permissions
('categories', 'View Categories', 'categories.view', 'Can view item categories', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('categories', 'Create Category', 'categories.create', 'Can create item category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('categories', 'Edit Category', 'categories.edit', 'Can edit item category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('categories', 'Delete Category', 'categories.delete', 'Can delete item category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

('sub_categories', 'View Sub-Categories', 'sub_categories.view', 'Can view item sub-categories', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sub_categories', 'Create Sub-Category', 'sub_categories.create', 'Can create item sub-category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sub_categories', 'Edit Sub-Category', 'sub_categories.edit', 'Can edit item sub-category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sub_categories', 'Delete Sub-Category', 'sub_categories.delete', 'Can delete item sub-category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= TAGS =================
('tags', 'View Tags', 'tags.view', 'Can view tags list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tags', 'Create Tag', 'tags.create', 'Can create new tag', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tags', 'Edit Tag', 'tags.edit', 'Can edit tag', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('tags', 'Delete Tag', 'tags.delete', 'Can delete tag', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= ITEM TYPES =================
('item_types', 'View Item Types', 'item_types.view', 'Can view item types list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('item_types', 'Create Item Type', 'item_types.create', 'Can create new item type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('item_types', 'Edit Item Type', 'item_types.edit', 'Can edit item type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('item_types', 'Delete Item Type', 'item_types.delete', 'Can delete item type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= MINOR CATEGORIES =================
('minor_categories', 'View Minor Categories', 'minor_categories.view', 'Can view minor categories list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('minor_categories', 'Create Minor Category', 'minor_categories.create', 'Can create new minor category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('minor_categories', 'Edit Minor Category', 'minor_categories.edit', 'Can edit minor category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('minor_categories', 'Delete Minor Category', 'minor_categories.delete', 'Can delete minor category', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= SALES SETUPS =================
('sales_setups', 'View Sales Setups', 'sales_setups.view', 'Can view sales setups list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_setups', 'Create Sales Setup', 'sales_setups.create', 'Can create new sales setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_setups', 'Edit Sales Setup', 'sales_setups.edit', 'Can edit sales setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_setups', 'Delete Sales Setup', 'sales_setups.delete', 'Can delete sales setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= CALCULATED BASE PRICES =================
('calculated_base_prices', 'View Calculated Base Prices', 'calculated_base_prices.view', 'Can view calculated base prices list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('calculated_base_prices', 'Create Calculated Base Price', 'calculated_base_prices.create', 'Can create new calculated base price', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('calculated_base_prices', 'Edit Calculated Base Price', 'calculated_base_prices.edit', 'Can edit calculated base price', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('calculated_base_prices', 'Delete Calculated Base Price', 'calculated_base_prices.delete', 'Can delete calculated base price', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= SALES SUPPLY TYPES =================
('sales_supply_types', 'View Sales Supply Types', 'sales_supply_types.view', 'Can view sales supply types list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_supply_types', 'Create Sales Supply Type', 'sales_supply_types.create', 'Can create new sales supply type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_supply_types', 'Edit Sales Supply Type', 'sales_supply_types.edit', 'Can edit sales supply type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_supply_types', 'Delete Sales Supply Type', 'sales_supply_types.delete', 'Can delete sales supply type', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= SALES TAX SETUPS =================
('sales_tax_setups', 'View Sales Tax Setups', 'sales_tax_setups.view', 'Can view sales tax setups list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_tax_setups', 'Create Sales Tax Setup', 'sales_tax_setups.create', 'Can create new sales tax setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_tax_setups', 'Edit Sales Tax Setup', 'sales_tax_setups.edit', 'Can edit sales tax setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('sales_tax_setups', 'Delete Sales Tax Setup', 'sales_tax_setups.delete', 'Can delete sales tax setup', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ================= ITEMS =================
('items', 'View Items', 'items.view', 'Can view items list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('items', 'Create Item', 'items.create', 'Can create new item', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('items', 'Edit Item', 'items.edit', 'Can edit item', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('items', 'Delete Item', 'items.delete', 'Can delete item', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP); 
---     umo permissions
('uom', 'View Units of Measurement', 'uom.view', 'Can view units of measurement list', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('uom', 'Create Unit of Measurement', 'uom.create', 'Can create new unit of measurement', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('uom', 'Edit Unit of Measurement', 'uom.edit', 'Can edit unit of measurement', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('uom', 'Delete Unit of Measurement', 'uom.delete', 'Can delete unit of measurement', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


