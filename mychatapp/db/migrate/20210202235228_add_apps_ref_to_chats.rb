class AddAppsRefToChats < ActiveRecord::Migration[6.1]
  def change
    add_reference :chats, :apps, null: false, foreign_key: true
  end
end
