class CreateMessages < ActiveRecord::Migration[6.1]
  def change
    create_table :messages do |t|
      t.integer :message_number
      t.string :message

      t.timestamps
    end
  end
end
