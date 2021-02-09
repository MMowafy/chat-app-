class AddChatsRefToMessages < ActiveRecord::Migration[6.1]
  def change
    add_reference :messages, :chats, null: false, foreign_key: true
  end
end
