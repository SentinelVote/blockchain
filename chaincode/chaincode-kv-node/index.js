const { Contract } = require("fabric-contract-api");
const crypto = require("crypto");

/**
 * Contract for key-value storage operations.
 */
class KVContract extends Contract {
  /**
   * Constructor for KVContract.
   */
  constructor() {
    super("KVContract");
  }

  /**
   * Function to be invoked on chaincode instantiation.
   * @param {Context} ctx - The transaction context.
   * @returns {Promise<void>} - A promise indicating successful instantiation.
   */
  async instantiate(ctx) {
    // Function that will be invoked on chaincode instantiation.
  }

  /**
   * Puts a key-value pair in the ledger.
   * @param {Context} ctx - The transaction context.
   * @param {string} key - The key to put in the ledger.
   * @param {string} value - The value to associate with the key.
   * @returns {Promise<{success: string}>} - The result of the operation.
   */
  async put(ctx, key, value) {
    await ctx.stub.putState(key, Buffer.from(value));
    return { success: "OK" };
  }

  /**
   * Retrieves a value from the ledger by its key.
   * @param {Context} ctx - The transaction context.
   * @param {string} key - The key of the value to retrieve.
   * @returns {Promise<{success: string} | {error: string}>} - The result of the operation.
   */
  async get(ctx, key) {
    const buffer = await ctx.stub.getState(key);
    if (!buffer || !buffer.length) return { error: "NOT_FOUND" };
    return { success: buffer.toString() };
  }

  /**
   * Stores a private message in a specified collection.
   * @param {Context} ctx - The transaction context.
   * @param {string} collection - The collection in which to store the message.
   * @returns {Promise<{success: string}>} - The result of the operation.
   */
  async putPrivateMessage(ctx, collection) {
    const transient = ctx.stub.getTransient();
    const message = transient.get("message");
    await ctx.stub.putPrivateData(collection, "message", message);
    return { success: "OK" };
  }

  /**
   * Retrieves a private message from a specified collection.
   * @param {Context} ctx - The transaction context.
   * @param {string} collection - The collection from which to retrieve the message.
   * @returns {Promise<{success: string}>} - The result of the operation.
   */
  async getPrivateMessage(ctx, collection) {
    const message = await ctx.stub.getPrivateData(collection, "message");
    const messageString = message.toBuffer ? message.toBuffer().toString() : message.toString();
    return { success: messageString };
  }

  /**
   * Verifies the hash of a private message against the stored hash in a specified collection.
   * @param {Context} ctx - The transaction context.
   * @param {string} collection - The collection in which the message is stored.
   * @returns {Promise<{success: string} | {error: string}>} - The result of the operation.
   */
  async verifyPrivateMessage(ctx, collection) {
    const transient = ctx.stub.getTransient();
    const message = transient.get("message");
    const messageString = message.toBuffer ? message.toBuffer().toString() : message.toString();
    const currentHash = crypto.createHash("sha256").update(messageString).digest("hex");
    const privateDataHash = (await ctx.stub.getPrivateDataHash(collection, "message")).toString("hex");
    if (privateDataHash !== currentHash) {
      return { error: "VERIFICATION_FAILED" };
    }
    return { success: "OK" };
  }
}

exports.contracts = [KVContract];
