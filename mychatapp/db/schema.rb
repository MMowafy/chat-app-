# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2021_02_02_235259) do

  create_table "apps", charset: "utf8mb4", force: :cascade do |t|
    t.string "name"
    t.string "token"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
  end

  create_table "chats", charset: "utf8mb4", force: :cascade do |t|
    t.integer "chat_number"
    t.integer "messages_count"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.bigint "apps_id", null: false
    t.index ["apps_id"], name: "index_chats_on_apps_id"
  end

  create_table "messages", charset: "utf8mb4", force: :cascade do |t|
    t.integer "message_number"
    t.string "message"
    t.datetime "created_at", precision: 6, null: false
    t.datetime "updated_at", precision: 6, null: false
    t.bigint "chats_id", null: false
    t.index ["chats_id"], name: "index_messages_on_chats_id"
  end

  add_foreign_key "chats", "apps", column: "apps_id"
  add_foreign_key "messages", "chats", column: "chats_id"
end
